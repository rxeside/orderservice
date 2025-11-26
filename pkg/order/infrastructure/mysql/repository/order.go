package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"gitea.xscloud.ru/xscloud/golib/pkg/infrastructure/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"orderservice/pkg/order/domain/model"
)

func NewOrderRepository(ctx context.Context, client mysql.ClientContext) model.OrderRepository {
	return &orderRepository{
		ctx:    ctx,
		client: client,
	}
}

type orderRepository struct {
	ctx    context.Context
	client mysql.ClientContext
}

func (r *orderRepository) NextID() (uuid.UUID, error) {
	return uuid.NewV7()
}

func (r *orderRepository) Store(order model.Order) error {
	_, err := r.client.ExecContext(r.ctx,
		`
	INSERT INTO orders (order_id, user_id, product_id, price, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
		status=VALUES(status),
	    updated_at=VALUES(updated_at)
	`,
		order.OrderID,
		order.UserID,
		order.ProductID,
		order.Price,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)
	return errors.WithStack(err)
}

func (r *orderRepository) Find(spec model.FindSpec) (*model.Order, error) {
	dto := struct {
		OrderID   uuid.UUID `db:"order_id"`
		UserID    uuid.UUID `db:"user_id"`
		ProductID uuid.UUID `db:"product_id"`
		Price     int64     `db:"price"`
		Status    int       `db:"status"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}{}
	query, args := r.buildSpecArgs(spec)

	err := r.client.GetContext(
		r.ctx,
		&dto,
		`SELECT order_id, user_id, product_id, price, status, created_at, updated_at FROM orders WHERE `+query,
		args...,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.WithStack(model.ErrOrderNotFound)
		}
		return nil, errors.WithStack(err)
	}

	return &model.Order{
		OrderID:   dto.OrderID,
		UserID:    dto.UserID,
		ProductID: dto.ProductID,
		Price:     dto.Price,
		Status:    model.OrderStatus(dto.Status),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}, nil
}

func (r *orderRepository) buildSpecArgs(spec model.FindSpec) (query string, args []interface{}) {
	var parts []string
	if spec.OrderID != nil {
		parts = append(parts, "order_id = ?")
		args = append(args, *spec.OrderID)
	}
	return strings.Join(parts, " AND "), args
}
