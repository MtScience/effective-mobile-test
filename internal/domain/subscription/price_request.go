package subscription

import (
	"time"

	"github.com/google/uuid"
)

type PriceRequest struct {
	UserID    *uuid.UUID `json:"user_id"`
	Service   *string    `json:"service"`
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
}
