package gateway

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"orderservice/api/clients/productinternal"
)

type ProductGateway interface {
	FindProduct(ctx context.Context, productID uuid.UUID) (name string, price int64, err error)
}

func NewProductGateway(conn *grpc.ClientConn) ProductGateway {
	return &productGateway{
		client: productinternal.NewProductInternalServiceClient(conn),
	}
}

type productGateway struct {
	client productinternal.ProductInternalServiceClient
}

func (g *productGateway) FindProduct(ctx context.Context, productID uuid.UUID) (string, int64, error) {
	resp, err := g.client.FindProduct(ctx, &productinternal.FindProductRequest{
		ProductID: productID.String(),
	})
	if err != nil {
		return "", 0, errors.WithStack(err)
	}
	if resp.Product == nil {
		return "", 0, errors.Errorf("product %s not found", productID)
	}
	return resp.Product.Name, resp.Product.Price, nil
}
