package handlers

import (
	domain "subscriptions/internal/domain/subscription"
	"subscriptions/internal/types"

	"github.com/google/uuid"
)

type subscriptionModificationDTO struct {
	UserID         uuid.UUID   `json:"user_id"`
	Service        string      `json:"service"`
	Price          int64       `json:"price"`
	SubscribedOn   types.Date  `json:"subscribed_on"`
	UnsubscribedOn *types.Date `json:"unsubscribed_on"`
}

type subscriptionDTO struct {
	ID             int64       `json:"id"`
	UserID         uuid.UUID   `json:"user_id"`
	Service        string      `json:"service"`
	Price          int64       `json:"price"`
	SubscribedOn   types.Date  `json:"subscribed_on"`
	UnsubscribedOn *types.Date `json:"unsubscribed_on"`
}

type priceRequestDTO struct {
	UserID    *uuid.UUID `json:"user_id"`
	Service   *string    `json:"service"`
	StartDate types.Date `json:"start_date"`
	EndDate   types.Date `json:"end_date"`
}

type priceDTO struct {
	Price int64 `json:"price"`
}

func newSubscriptionDTO(sub *domain.Subscription) (dto subscriptionDTO) {
	dto = subscriptionDTO{
		ID:           sub.ID,
		UserID:       sub.UserID,
		Service:      sub.Service,
		Price:        sub.Price,
		SubscribedOn: types.Date{sub.SubscribedOn},
	}

	if sub.UnsubscribedOn != nil {
		dto.UnsubscribedOn = &types.Date{*sub.UnsubscribedOn}
	} else {
		dto.UnsubscribedOn = nil
	}

	return
}
