package domain

import (
	"github.com/google/uuid"
)

type PaymentPayload struct {
	UserID     uuid.UUID
	PaymentSum int64
}
