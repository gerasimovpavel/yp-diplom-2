package repoPostgres

import (
	"context"
	"errors"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"yp-diplom-2/cmd/internal/config"
	"yp-diplom-2/cmd/internal/repository"
)

func NewRepositories(cfg *config.Config) (*repository.Repositories, error) {
	ctx := context.Background()
	URI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.URL, cfg.Database.Port, cfg.Database.Db,
	)
	pool, err := pgxpool.New(ctx, URI)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &repository.Repositories{
		Users: NewUserRepo(pool, trmpgx.DefaultCtxGetter),
	}, nil
}

func IsDuplicate(err error) bool {
	var e *pgconn.PgError
	if errors.As(err, &e) {
		if e.Code == pgerrcode.UniqueViolation {
			return true
		}
	}
	return false
}
