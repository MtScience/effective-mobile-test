package subscription

import (
	"context"
	"fmt"
	"strings"
	domain "subscriptions/internal/domain/subscription"
	"subscriptions/internal/logging"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*domain.Subscription, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &domain.Subscription{
		UserID:       normalized.UserID,
		Service:      normalized.Service,
		Price:        normalized.Price,
		SubscribedOn: normalized.SubscribedOn.Time,
	}
	if normalized.UnsubscribedOn != nil {
		model.UnsubscribedOn = &normalized.UnsubscribedOn.Time
	} else {
		model.UnsubscribedOn = nil
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*domain.Subscription, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*domain.Subscription, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &domain.Subscription{
		ID:           id,
		UserID:       normalized.UserID,
		Service:      normalized.Service,
		Price:        normalized.Price,
		SubscribedOn: normalized.SubscribedOn.Time,
	}
	if normalized.UnsubscribedOn != nil {
		model.UnsubscribedOn = &normalized.UnsubscribedOn.Time
	} else {
		model.UnsubscribedOn = nil
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]domain.Subscription, error) {
	return s.repo.List(ctx)
}

func (s *Service) TotalPrice(ctx context.Context, priceRequest PriceRequestInput) (int64, error) {
	normalized, err := validatePriceRequestInput(priceRequest)
	if err != nil {
		return -1, err
	}

	return s.repo.TotalPrice(
		ctx, &domain.PriceRequest{
			UserID:    normalized.UserID,
			Service:   normalized.Service,
			StartDate: normalized.StartDate.Time,
			EndDate:   normalized.EndDate.Time,
		},
	)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	logging.Logger.Info("validating input for subscription creation")

	input.Service = strings.TrimSpace(input.Service)

	if input.Service == "" {
		return CreateInput{}, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}

	if input.Price < 0 {
		return CreateInput{}, fmt.Errorf("%w: price can't be negative", ErrInvalidInput)
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	logging.Logger.Info("validating input for subscription update")

	input.Service = strings.TrimSpace(input.Service)

	if input.Service == "" {
		return UpdateInput{}, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}

	if input.Price < 0 {
		return UpdateInput{}, fmt.Errorf("%w: price can't be negative", ErrInvalidInput)
	}

	return input, nil
}

func validatePriceRequestInput(input PriceRequestInput) (PriceRequestInput, error) {
	logging.Logger.Info("validating input for subscription price request")

	if input.Service != nil {
		service := strings.TrimSpace(*input.Service)
		if service == "" {
			input.Service = nil
		} else {
			input.Service = &service
		}
	}

	if input.UserID == nil && input.Service == nil {
		return PriceRequestInput{}, fmt.Errorf("%w: user_id and service can't both be null", ErrInvalidInput)
	}

	return input, nil
}
