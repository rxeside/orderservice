package transport

import (
	"context"

	"github.com/google/uuid"

	"orderservice/api/server/orderinternal"
	"orderservice/pkg/order/application/service"
)

func NewOrderInternalAPI(orderService service.OrderService) orderinternal.OrderInternalServiceServer {
	return &orderInternalAPI{orderService: orderService}
}

type orderInternalAPI struct {
	orderService service.OrderService
	orderinternal.UnimplementedOrderInternalServiceServer
}

func (a *orderInternalAPI) CreateOrder(ctx context.Context, req *orderinternal.CreateOrderRequest) (*orderinternal.CreateOrderResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, err
	}

	orderID, err := a.orderService.CreateOrder(ctx, userID, productID)
	if err != nil {
		return nil, err
	}
	return &orderinternal.CreateOrderResponse{OrderID: orderID.String()}, nil
}
func (a *orderInternalAPI) FindOrder(ctx context.Context, req *orderinternal.FindOrderRequest) (*orderinternal.FindOrderResponse, error) {
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, err
	}
	order, err := a.orderService.FindOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &orderinternal.FindOrderResponse{
		Order: &orderinternal.Order{
			OrderID:   order.OrderID.String(),
			UserID:    order.UserID.String(),
			ProductID: order.ProductID.String(),
			Price:     order.Price,
			Status:    orderinternal.OrderStatus(order.Status),
		},
	}, nil
}
