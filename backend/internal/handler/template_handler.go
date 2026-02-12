package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"reelcut/internal/middleware"
	"reelcut/internal/service"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TemplateHandler struct {
	templateSvc *service.TemplateService
}

func NewTemplateHandler(templateSvc *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{templateSvc: templateSvc}
}

// List godoc
// @Summary		List current user's templates
// @Tags			templates
// @Produce		json
// @Security	BearerAuth
// @Param		page		query		int	false	"Page"
// @Param		per_page	query		int	false	"Per page"
// @Success	200	{object}	object
// @Router		/api/v1/templates [get]
func (h *TemplateHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	uid := &userID
	list, total, err := h.templateSvc.List(c.Request.Context(), uid, false, page, perPage)
	if err != nil {
		utils.Internal(c, "")
		return
	}
	utils.JSONPaginated(c, http.StatusOK, gin.H{"templates": list}, page, perPage, total)
}

// GetPublicTemplates godoc
// @Summary		List public templates (no auth)
// @Tags			templates
// @Produce		json
// @Param		page		query		int	false	"Page"
// @Param		per_page	query		int	false	"Per page"
// @Success	200	{object}	object
// @Router		/api/v1/templates/public [get]
func (h *TemplateHandler) GetPublicTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	list, total, err := h.templateSvc.List(c.Request.Context(), nil, true, page, perPage)
	if err != nil {
		utils.Internal(c, "")
		return
	}
	utils.JSONPaginated(c, http.StatusOK, gin.H{"templates": list}, page, perPage, total)
}

// Create godoc
// @Summary		Create a template
// @Tags			templates
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		body	body		object	true	"name, category, is_public, style_config"
// @Success	201	{object}	object
// @Router		/api/v1/templates [post]
func (h *TemplateHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	uid, _ := uuid.Parse(userID)
	var body struct {
		Name        string          `json:"name" binding:"required"`
		Category    string          `json:"category"`
		IsPublic    bool            `json:"is_public"`
		StyleConfig json.RawMessage `json:"style_config" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, nil)
		return
	}
	t, err := h.templateSvc.Create(c.Request.Context(), &uid, body.Name, body.Category, body.IsPublic, body.StyleConfig)
	if err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"template": t})
}

// GetByID godoc
// @Summary		Get template by ID
// @Tags			templates
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Template ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/templates/{id} [get]
func (h *TemplateHandler) GetByID(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	t, err := h.templateSvc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil || t == nil {
		utils.NotFound(c, "Template not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"template": t})
}

// Update godoc
// @Summary		Update a template
// @Tags			templates
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id		path		string	true	"Template ID"
// @Param		body	body		object	true	"name, category, is_public, style_config"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/templates/{id} [put]
func (h *TemplateHandler) Update(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	t, err := h.templateSvc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil || t == nil {
		utils.NotFound(c, "Template not found")
		return
	}
	var body struct {
		Name        *string         `json:"name"`
		Category    *string         `json:"category"`
		IsPublic    *bool           `json:"is_public"`
		StyleConfig json.RawMessage `json:"style_config"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, nil)
		return
	}
	if body.Name != nil {
		t.Name = *body.Name
	}
	if body.Category != nil {
		t.Category = body.Category
	}
	if body.IsPublic != nil {
		t.IsPublic = *body.IsPublic
	}
	if body.StyleConfig != nil {
		t.StyleConfig = body.StyleConfig
	}
	if err := h.templateSvc.Update(c.Request.Context(), t); err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"template": t})
}

// Delete godoc
// @Summary		Delete a template
// @Tags			templates
// @Security	BearerAuth
// @Param		id	path		string	true	"Template ID"
// @Success	204	"No Content"
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/templates/{id} [delete]
func (h *TemplateHandler) Delete(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	if err := h.templateSvc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		utils.NotFound(c, "Template not found")
		return
	}
	c.Status(http.StatusNoContent)
}
