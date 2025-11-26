package service

import (
	"time"

	"github.com/google/uuid"

	"orderservice/pkg/common/domain"
	"orderservice/pkg/order/domain/model"
)

type OrderService interface {
	CreateOrder(userID, productID uuid.UUID, price int64) (uuid.UUID, error)
}

func NewOrderService(repo model.OrderRepository, dispatcher domain.EventDispatcher) OrderService {
	return &orderService{
		repo:       repo,
		dispatcher: dispatcher,
	}
}

type orderService struct {
	repo       model.OrderRepository
	dispatcher domain.EventDispatcher
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

	// отправляем событие для хореографии
	return id, s.dispatcher.Dispatch(model.OrderCreated{
		OrderID:   id,
		UserID:    userID,
		ProductID: productID,
		Price:     price,
		CreatedAt: now,
	})
}
