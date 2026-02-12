package handler

import (
	"net/http"
	"strconv"

	"reelcut/internal/middleware"
	"reelcut/internal/repository"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	jobRepo repository.ProcessingJobRepository
}

func NewJobHandler(jobRepo repository.ProcessingJobRepository) *JobHandler {
	return &JobHandler{jobRepo: jobRepo}
}

// List godoc
// @Summary		List processing jobs
// @Tags			jobs
// @Produce		json
// @Security	BearerAuth
// @Param		page		query		int		false	"Page"
// @Param		per_page	query		int		false	"Per page"
// @Param		status		query		string	false	"Filter by status"
// @Success	200	{object}	object
// @Router		/api/v1/jobs [get]
func (h *JobHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	status := c.Query("status")
	var st *string
	if status != "" {
		st = &status
	}
	list, total, err := h.jobRepo.ListByUserID(c.Request.Context(), userID, st, perPage, (page-1)*perPage)
	if err != nil {
		utils.Internal(c, "")
		return
	}
	utils.JSONPaginated(c, http.StatusOK, gin.H{"jobs": list}, page, perPage, total)
}

// GetByID godoc
// @Summary		Get job by ID
// @Tags			jobs
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Job ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/jobs/{id} [get]
func (h *JobHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	j, err := h.jobRepo.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil || j == nil || j.UserID.String() != userID {
		utils.NotFound(c, "Job not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"job": j})
}

// Cancel godoc
// @Summary		Cancel a job
// @Tags			jobs
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Job ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/jobs/{id}/cancel [post]
func (h *JobHandler) Cancel(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	j, err := h.jobRepo.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil || j == nil || j.UserID.String() != userID {
		utils.NotFound(c, "Job not found")
		return
	}
	j.Status = "cancelled"
	if err := h.jobRepo.Update(c.Request.Context(), j); err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job cancelled"})
}
