package service

import (
	"context"

	"gitea.xscloud.ru/xscloud/golib/pkg/application/outbox"
	"github.com/google/uuid"

	"orderservice/pkg/common/domain"
	"orderservice/pkg/order/domain/model"
	"orderservice/pkg/order/domain/service"
	"orderservice/pkg/order/infrastructure/transport/gateway"
)

type RepositoryProvider interface {
	OrderRepository(ctx context.Context) model.OrderRepository
}

type LockableUnitOfWork interface {
	Execute(ctx context.Context, lockNames []string, f func(provider RepositoryProvider) error) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, userID, productID uuid.UUID) (uuid.UUID, error)
	FindOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
}

func NewOrderService(
	luow LockableUnitOfWork,
	userGateway gateway.UserGateway,
	productGateway gateway.ProductGateway,
	eventDispatcher outbox.EventDispatcher[outbox.Event],
) OrderService {
	return &orderServiceImpl{
		luow:            luow,
		userGateway:     userGateway,
		productGateway:  productGateway,
		eventDispatcher: eventDispatcher,
	}
}

type orderServiceImpl struct {
	luow            LockableUnitOfWork
	userGateway     gateway.UserGateway
	productGateway  gateway.ProductGateway
	eventDispatcher outbox.EventDispatcher[outbox.Event]
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, userID, productID uuid.UUID) (uuid.UUID, error) {
	// оркестрация: проверяем пользователя
	if err := s.userGateway.CheckUserActive(ctx, userID); err != nil {
		return uuid.Nil, err
	}

	// оркестрация: получаем актуальную цену продукта
	_, price, err := s.productGateway.FindProduct(ctx, productID)
	if err != nil {
		return uuid.Nil, err
	}

	var orderID uuid.UUID
	err = s.luow.Execute(ctx, nil, func(provider RepositoryProvider) error {
		domainDispatcher := &domainEventDispatcher{
			ctx:             ctx,
			eventDispatcher: s.eventDispatcher,
		}

		domainSvc := service.NewOrderService(provider.OrderRepository(ctx), domainDispatcher)

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

type domainEventDispatcher struct {
	ctx             context.Context
	eventDispatcher outbox.EventDispatcher[outbox.Event]
}

func (d *domainEventDispatcher) Dispatch(event domain.Event) error {
	return d.eventDispatcher.Dispatch(d.ctx, event)
}
