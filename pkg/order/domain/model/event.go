package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderCreated struct {
	OrderID   uuid.UUID
	UserID    uuid.UUID
	ProductID uuid.UUID
	Price     int64
	CreatedAt time.Time
}

func (e OrderCreated) Type() string {
	return "order_created"
}
