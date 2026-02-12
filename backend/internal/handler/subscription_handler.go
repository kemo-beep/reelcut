package handler

import (
	"context"
	"net/http"

	"reelcut/internal/domain"
	"reelcut/internal/middleware"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subscriptionSvc SubscriptionService
}

type SubscriptionService interface {
	GetByUserID(ctx context.Context, userID string) (*domain.Subscription, error)
	CreateSubscription(ctx context.Context, userID, tier string) (*domain.Subscription, error)
	CancelSubscription(ctx context.Context, userID string) error
	UpdateSubscription(ctx context.Context, userID, tier string) error
}

func NewSubscriptionHandler(subscriptionSvc SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionSvc: subscriptionSvc}
}

// GetMySubscription returns the current user's subscription.
// @Summary	Get my subscription
// @Tags		subscription
// @Produce	json
// @Security	BearerAuth
// @Success	200	{object}	domain.Subscription
// @Failure	404	{object}	utils.ErrorResponse
// @Router	/users/me/subscription [get]
func (h *SubscriptionHandler) GetMySubscription(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	sub, err := h.subscriptionSvc.GetByUserID(c.Request.Context(), userID)
	if err != nil || sub == nil {
		utils.Error(c, http.StatusNotFound, "NOT_FOUND", "subscription not found", nil)
		return
	}
	c.JSON(http.StatusOK, sub)
}

// CreateSubscription creates a new subscription for the current user.
// @Summary	Create subscription
// @Tags		subscription
// @Accept	json
// @Produce	json
// @Security	BearerAuth
// @Param	body	body		object	true	"tier"
// @Success	200	{object}	domain.Subscription
// @Failure	400	{object}	utils.ErrorResponse
// @Router	/subscriptions/create [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var body struct {
		Tier string `json:"tier"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error(), nil)
		return
	}
	if body.Tier == "" {
		body.Tier = "pro"
	}
	sub, err := h.subscriptionSvc.CreateSubscription(c.Request.Context(), userID, body.Tier)
	if err != nil {
		if err == domain.ErrNotFound {
			utils.Error(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
			return
		}
		utils.Error(c, http.StatusBadRequest, "SUBSCRIPTION_ERROR", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, sub)
}

// CancelSubscription cancels the current user's subscription.
// @Summary	Cancel subscription
// @Tags		subscription
// @Security	BearerAuth
// @Success	204
// @Failure	404	{object}	utils.ErrorResponse
// @Router	/subscriptions/cancel [post]
func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err := h.subscriptionSvc.CancelSubscription(c.Request.Context(), userID); err != nil {
		if err == domain.ErrNotFound {
			utils.Error(c, http.StatusNotFound, "NOT_FOUND", "subscription not found", nil)
			return
		}
		utils.Error(c, http.StatusBadRequest, "CANCEL_ERROR", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateSubscription updates the current user's subscription tier.
// @Summary	Update subscription
// @Tags		subscription
// @Accept	json
// @Security	BearerAuth
// @Param	body	body		object	true	"tier"
// @Success	204
// @Failure	400	{object}	utils.ErrorResponse
// @Router	/subscriptions/update [post]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var body struct {
		Tier string `json:"tier"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error(), nil)
		return
	}
	if err := h.subscriptionSvc.UpdateSubscription(c.Request.Context(), userID, body.Tier); err != nil {
		if err == domain.ErrNotFound {
			utils.Error(c, http.StatusNotFound, "NOT_FOUND", "subscription not found", nil)
			return
		}
		utils.Error(c, http.StatusBadRequest, "UPDATE_ERROR", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}
