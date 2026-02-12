package service

import (
	"context"
	"fmt"
	"time"

	"reelcut/internal/domain"
	"reelcut/internal/repository"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/subscription"
)

type SubscriptionService struct {
	subRepo  repository.SubscriptionRepository
	userRepo repository.UserRepository
	secretKey string
	priceIDPro string
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository, userRepo repository.UserRepository, secretKey, priceIDPro string) *SubscriptionService {
	return &SubscriptionService{subRepo: subRepo, userRepo: userRepo, secretKey: secretKey, priceIDPro: priceIDPro}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, userID, tier string) (*domain.Subscription, error) {
	if s.secretKey == "" {
		return nil, fmt.Errorf("stripe not configured")
	}
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || u == nil {
		return nil, domain.ErrNotFound
	}
	stripe.Key = s.secretKey

	// Create or reuse Stripe customer (by email for idempotency)
	custParams := &stripe.CustomerParams{
		Email: stripe.String(u.Email),
	}
	cust, err := customer.New(custParams)
	if err != nil {
		return nil, fmt.Errorf("stripe customer: %w", err)
	}

	priceID := s.priceIDPro
	subParams := &stripe.SubscriptionParams{
		Customer: stripe.String(cust.ID),
		Items: []*stripe.SubscriptionItemsParams{
			{Price: stripe.String(priceID)},
		},
	}
	sub, err := subscription.New(subParams)
	if err != nil {
		return nil, fmt.Errorf("stripe subscription: %w", err)
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

	dom := &domain.Subscription{
		ID:                   uuid.New(),
		UserID:               uuid.MustParse(userID),
		Tier:                 tier,
		Status:               string(sub.Status),
		StripeSubscriptionID: &sub.ID,
		CurrentPeriodStart:   periodStart,
		CurrentPeriodEnd:     periodEnd,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	if err := s.subRepo.Create(ctx, dom); err != nil {
		return nil, err
	}
	// Optionally update user tier and credits
	u.SubscriptionTier = tier
	_ = s.userRepo.Update(ctx, u)
	return dom, nil
}

func (s *SubscriptionService) CancelSubscription(ctx context.Context, userID string) error {
	sub, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil || sub == nil || sub.StripeSubscriptionID == nil {
		return domain.ErrNotFound
	}
	stripe.Key = s.secretKey
	_, err = subscription.Cancel(*sub.StripeSubscriptionID, nil)
	if err != nil {
		return fmt.Errorf("stripe cancel: %w", err)
	}
	sub.Status = "cancelled"
	return s.subRepo.Update(ctx, sub)
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, userID, tier string) error {
	sub, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil || sub == nil || sub.StripeSubscriptionID == nil {
		return domain.ErrNotFound
	}
	stripe.Key = s.secretKey
	// Update subscription items to new price
	subStripe, err := subscription.Get(*sub.StripeSubscriptionID, nil)
	if err != nil || subStripe == nil || len(subStripe.Items.Data) == 0 {
		return fmt.Errorf("no subscription items")
	}
	itemID := subStripe.Items.Data[0].ID
	_, err = subscription.Update(*sub.StripeSubscriptionID, &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{ID: stripe.String(itemID), Price: stripe.String(s.priceIDPro)},
		},
	})
	if err != nil {
		return fmt.Errorf("stripe update: %w", err)
	}
	sub.Tier = tier
	return s.subRepo.Update(ctx, sub)
}

func (s *SubscriptionService) GetByUserID(ctx context.Context, userID string) (*domain.Subscription, error) {
	return s.subRepo.GetByUserID(ctx, userID)
}
