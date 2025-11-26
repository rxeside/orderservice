package service

import (
	"context"

	"github.com/google/uuid"

	"orderservice/pkg/order/domain/model"
	"orderservice/pkg/order/domain/service"
)

type RepositoryProvider interface {
	OrderRepository(ctx context.Context) model.OrderRepository
}

type LockableUnitOfWork interface {
	Execute(ctx context.Context, lockNames []string, f func(provider RepositoryProvider) error) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, userID, productID uuid.UUID, price int64) (uuid.UUID, error)
	FindOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
}

func NewOrderService(luow LockableUnitOfWork) OrderService {
	return &orderServiceImpl{luow: luow}
}

type orderServiceImpl struct {
	luow LockableUnitOfWork
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, userID, productID uuid.UUID, price int64) (uuid.UUID, error) {
	var orderID uuid.UUID
	// Пока лок не нужен, но структура готова
	err := s.luow.Execute(ctx, nil, func(provider RepositoryProvider) error {
		domainSvc := service.NewOrderService(provider.OrderRepository(ctx))
		var err error
		orderID, err = domainSvc.CreateOrder(userID, productID, price)
		return err
	})
	return orderID, err
}

func (s *orderServiceImpl) FindOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	var order *model.Order
	err := s.luow.Execute(ctx, nil, func(provider RepositoryProvider) error {
		repo := provider.OrderRepository(ctx)
		var err error
		order, err = repo.Find(model.FindSpec{OrderID: &orderID})
		return err
	})
	return order, err
}
