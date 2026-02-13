package handler

import (
	"net/http"
	"strconv"

	"reelcut/internal/domain"
	"reelcut/internal/middleware"
	"reelcut/internal/service"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BrollHandler struct {
	clipSvc *service.ClipService
}

func NewBrollHandler(clipSvc *service.ClipService) *BrollHandler {
	return &BrollHandler{clipSvc: clipSvc}
}

// CreateAsset godoc
// @Summary		Upload a B-roll asset
// @Tags			broll
// @Accept		multipart/form-data
// @Produce		json
// @Security	BearerAuth
// @Param		file	formData	file	true	"Video file"
// @Param		project_id	formData	string	false	"Project ID"
// @Success	201	{object}	object
// @Router		/api/v1/broll/assets [post]
func (h *BrollHandler) CreateAsset(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		utils.ValidationError(c, []utils.ErrorDetail{{Field: "file", Message: "file is required"}})
		return
	}
	projectID := c.PostForm("project_id")
	var projID *string
	if projectID != "" {
		projID = &projectID
	}
	uid, _ := uuid.Parse(userID)
	f, err := file.Open()
	if err != nil {
		utils.Internal(c, "")
		return
	}
	defer f.Close()
	asset, err := h.clipSvc.CreateBrollAsset(c.Request.Context(), uid, projID, file.Filename, f, file.Header.Get("Content-Type"), file.Size)
	if err != nil {
		if err == domain.ErrValidation {
			utils.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "File too large or invalid", nil)
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"asset": asset})
}

// ListAssets godoc
// @Summary		List B-roll assets for the current user
// @Tags			broll
// @Produce		json
// @Security	BearerAuth
// @Param		project_id	query	string	false	"Filter by project"
// @Param		limit	query	int	false	"Limit"
// @Param		offset	query	int	false	"Offset"
// @Success	200	{object}	object
// @Router		/api/v1/broll/assets [get]
func (h *BrollHandler) ListAssets(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	projectID := c.Query("project_id")
	var projID *string
	if projectID != "" {
		projID = &projectID
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	list, total, err := h.clipSvc.ListBrollAssets(c.Request.Context(), userID, projID, limit, offset)
	if err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"assets": list, "total": total})
}
