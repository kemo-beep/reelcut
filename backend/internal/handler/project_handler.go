package handler

import (
	"net/http"
	"strconv"

	"reelcut/internal/domain"
	"reelcut/internal/middleware"
	"reelcut/internal/repository"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectRepo repository.ProjectRepository
}

func NewProjectHandler(projectRepo repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{projectRepo: projectRepo}
}

// List godoc
// @Summary		List projects
// @Tags			projects
// @Produce		json
// @Security	BearerAuth
// @Param		page		query		int	false	"Page"
// @Param		per_page	query		int	false	"Per page"
// @Success	200	{object}	object
// @Router		/api/v1/projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	list, total, err := h.projectRepo.ListByUserID(c.Request.Context(), userID, perPage, offset)
	if err != nil {
		utils.Internal(c, "")
		return
	}
	utils.JSONPaginated(c, http.StatusOK, gin.H{"projects": list}, page, perPage, total)
}

// Create godoc
// @Summary		Create a project
// @Tags			projects
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		body	body		object	true	"name, description"
// @Success	201	{object}	object
// @Router		/api/v1/projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	var body struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, []utils.ErrorDetail{{Message: err.Error()}})
		return
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.Unauthorized(c, "")
		return
	}
	p := &domain.Project{
		ID:          uuid.New(),
		UserID:      uid,
		Name:        body.Name,
		Description: body.Description,
	}
	if err := h.projectRepo.Create(c.Request.Context(), p); err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"project": p})
}

// GetByID godoc
// @Summary		Get project by ID
// @Tags			projects
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Project ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/projects/{id} [get]
func (h *ProjectHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	id := c.Param("id")
	p, err := h.projectRepo.GetByID(c.Request.Context(), id)
	if err != nil || p == nil {
		utils.NotFound(c, "Project not found")
		return
	}
	if p.UserID.String() != userID {
		utils.NotFound(c, "Project not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"project": p})
}

// Update godoc
// @Summary		Update a project
// @Tags			projects
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id		path		string	true	"Project ID"
// @Param		body	body		object	true	"name, description"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	id := c.Param("id")
	p, err := h.projectRepo.GetByID(c.Request.Context(), id)
	if err != nil || p == nil || p.UserID.String() != userID {
		utils.NotFound(c, "Project not found")
		return
	}
	var body struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, nil)
		return
	}
	if body.Name != nil {
		p.Name = *body.Name
	}
	if body.Description != nil {
		p.Description = body.Description
	}
	if err := h.projectRepo.Update(c.Request.Context(), p); err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"project": p})
}

// Delete godoc
// @Summary		Delete a project
// @Tags			projects
// @Security	BearerAuth
// @Param		id	path		string	true	"Project ID"
// @Success	204	"No Content"
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	id := c.Param("id")
	p, err := h.projectRepo.GetByID(c.Request.Context(), id)
	if err != nil || p == nil || p.UserID.String() != userID {
		utils.NotFound(c, "Project not found")
		return
	}
	if err := h.projectRepo.Delete(c.Request.Context(), id); err != nil {
		utils.Internal(c, "")
		return
	}
	c.Status(http.StatusNoContent)
}
