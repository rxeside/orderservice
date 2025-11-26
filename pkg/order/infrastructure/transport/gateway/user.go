package gateway

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"orderservice/api/clients/userinternal"
)

type UserGateway interface {
	CheckUserActive(ctx context.Context, userID uuid.UUID) error
}

func NewUserGateway(conn *grpc.ClientConn) UserGateway {
	return &userGateway{
		client: userinternal.NewUserInternalServiceClient(conn),
	}
}

type userGateway struct {
	client userinternal.UserInternalServiceClient
}

func (g *userGateway) CheckUserActive(ctx context.Context, userID uuid.UUID) error {
	resp, err := g.client.FindUser(ctx, &userinternal.FindUserRequest{
		UserID: userID.String(),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	if resp.User == nil {
		return errors.Errorf("user %s not found", userID)
	}

	if resp.User.Status != userinternal.UserStatus_Active {
		return errors.Errorf("user %s is not active (status: %v)", userID, resp.User.Status)
	}
	return nil
}
