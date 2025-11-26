package integrationevent

import (
	"encoding/json"

	"gitea.xscloud.ru/xscloud/golib/pkg/application/outbox"
	"github.com/pkg/errors"

	"orderservice/pkg/order/domain/model"
)

func NewEventSerializer() outbox.EventSerializer[outbox.Event] {
	return &eventSerializer{}
}

type eventSerializer struct{}

func (s eventSerializer) Serialize(event outbox.Event) (string, error) {
	switch e := event.(type) {
	case model.OrderCreated:
		b, err := json.Marshal(OrderCreated{
			OrderID:   e.OrderID.String(),
			UserID:    e.UserID.String(),
			ProductID: e.ProductID.String(),
			Price:     e.Price,
			CreatedAt: e.CreatedAt.Unix(),
		})
		return string(b), errors.WithStack(err)
	default:
		return "", errors.Errorf("unknown event %q", event.Type())
	}
}

type OrderCreated struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Price     int64  `json:"price"`
	CreatedAt int64  `json:"created_at"`
}
