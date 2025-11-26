package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderStatus int

const (
	Created OrderStatus = iota
	Paid
	Cancelled
	Completed
)

type Order struct {
	OrderID   uuid.UUID
	UserID    uuid.UUID
	ProductID uuid.UUID
	Price     int64
	Status    OrderStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FindSpec struct {
	OrderID *uuid.UUID
}

type OrderRepository interface {
	NextID() (uuid.UUID, error)
	Store(order Order) error
	Find(spec FindSpec) (*Order, error)
}
