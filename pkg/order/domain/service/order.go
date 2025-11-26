package service

import (
	"time"

	"github.com/google/uuid"

	"orderservice/pkg/order/domain/model"
)

type OrderService interface {
	CreateOrder(userID, productID uuid.UUID, price int64) (uuid.UUID, error)
	// Методы для обновления статуса добавим, когда будем реализовывать Saga
}

func NewOrderService(repo model.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

type orderService struct {
	repo model.OrderRepository
}

func (s *orderService) CreateOrder(userID, productID uuid.UUID, price int64) (uuid.UUID, error) {
	id, err := s.repo.NextID()
	if err != nil {
		return uuid.Nil, err
	}

	now := time.Now()
	order := model.Order{
		OrderID:   id,
		UserID:    userID,
		ProductID: productID,
		Price:     price,
		Status:    model.Created,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.repo.Store(order)
	if err != nil {
		return uuid.Nil, err
	}

	// Тут в будущем будет отправка события OrderCreated
	return id, nil
}
