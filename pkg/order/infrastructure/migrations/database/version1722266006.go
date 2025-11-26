package database

import (
	"context"

	"gitea.xscloud.ru/xscloud/golib/pkg/infrastructure/migrator"
	"gitea.xscloud.ru/xscloud/golib/pkg/infrastructure/mysql"
	"github.com/pkg/errors"
)

func NewVersion1722266006(client mysql.ClientContext) migrator.Migration {
	return &version1722266006{
		client: client,
	}
}

type version1722266006 struct {
	client mysql.ClientContext
}

func (v version1722266006) Version() int64 {
	return 1722266006
}

func (v version1722266006) Description() string {
	return "Create 'orders' table"
}

func (v version1722266006) Up(ctx context.Context) error {
	_, err := v.client.ExecContext(ctx, `
		CREATE TABLE orders
		(
			order_id      VARCHAR(64)  NOT NULL,
			user_id       VARCHAR(64)  NOT NULL,
			product_id    VARCHAR(64)  NOT NULL,
			price         BIGINT       NOT NULL,
			status        INT          NOT NULL,
			created_at    DATETIME     NOT NULL,
			updated_at    DATETIME     NOT NULL,
			PRIMARY KEY (order_id)
		)
			ENGINE = InnoDB
			CHARACTER SET = utf8mb4
			COLLATE utf8mb4_unicode_ci
	`)
	return errors.WithStack(err)
}
