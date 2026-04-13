package subscription

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID             int64      `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	Service        string     `json:"service"`
	Price          int64      `json:"price"`
	SubscribedOn   time.Time  `json:"subscribed_on"`
	UnsubscribedOn *time.Time `json:"unsubscribed_on"`
}
