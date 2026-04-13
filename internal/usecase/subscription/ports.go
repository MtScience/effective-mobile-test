package subscription

import (
	"context"
	domain "subscriptions/internal/domain/subscription"
	"subscriptions/internal/types"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	GetByID(ctx context.Context, id int64) (*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]domain.Subscription, error)
	TotalPrice(ctx context.Context, data *domain.PriceRequest) (int64, error)
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*domain.Subscription, error)
	GetByID(ctx context.Context, id int64) (*domain.Subscription, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*domain.Subscription, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]domain.Subscription, error)
	TotalPrice(ctx context.Context, priceRequest PriceRequestInput) (int64, error)
}

type CreateInput struct {
	UserID         uuid.UUID
	Service        string
	Price          int64
	SubscribedOn   types.Date
	UnsubscribedOn *types.Date
}

type UpdateInput struct {
	UserID         uuid.UUID
	Service        string
	Price          int64
	SubscribedOn   types.Date
	UnsubscribedOn *types.Date
}

type PriceRequestInput struct {
	UserID    *uuid.UUID
	Service   *string
	StartDate types.Date
	EndDate   types.Date
}
