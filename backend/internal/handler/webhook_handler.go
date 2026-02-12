package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"reelcut/internal/middleware"
	"reelcut/internal/repository"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

type WebhookHandler struct {
	stripeSecret string
	subRepo      repository.SubscriptionRepository
	userRepo     repository.UserRepository
}

func NewWebhookHandler(stripeWebhookSecret string, subRepo repository.SubscriptionRepository, userRepo repository.UserRepository) *WebhookHandler {
	return &WebhookHandler{stripeSecret: stripeWebhookSecret, subRepo: subRepo, userRepo: userRepo}
}

func (h *WebhookHandler) ProcessingComplete(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (h *WebhookHandler) Stripe(c *gin.Context) {
	if h.stripeSecret == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	signature := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(body, signature, h.stripeSecret)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	switch event.Type {
	case "customer.subscription.updated", "customer.subscription.created":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		h.upsertSubscriptionFromStripe(c.Request.Context(), &sub)
	case "customer.subscription.deleted":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		h.markSubscriptionCancelled(c.Request.Context(), sub.ID)
	}
	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *WebhookHandler) upsertSubscriptionFromStripe(ctx context.Context, sub *stripe.Subscription) {
	existing, _ := h.subRepo.GetByStripeID(ctx, sub.ID)
	if existing == nil {
		return
	}
	var periodStart, periodEnd *time.Time
	if sub.CurrentPeriodStart > 0 {
		t := time.Unix(sub.CurrentPeriodStart, 0)
		periodStart = &t
	}
	if sub.CurrentPeriodEnd > 0 {
		t := time.Unix(sub.CurrentPeriodEnd, 0)
		periodEnd = &t
	}
	existing.Status = string(sub.Status)
	existing.CurrentPeriodStart = periodStart
	existing.CurrentPeriodEnd = periodEnd
	_ = h.subRepo.Update(ctx, existing)
}

func (h *WebhookHandler) markSubscriptionCancelled(ctx context.Context, stripeSubID string) {
	existing, _ := h.subRepo.GetByStripeID(ctx, stripeSubID)
	if existing != nil {
		existing.Status = "cancelled"
		_ = h.subRepo.Update(ctx, existing)
	}
}
