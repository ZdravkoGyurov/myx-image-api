package postgresql

import (
	"context"

	"github.com/ZdravkoGyurov/myx-image-api/pkg/config"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
	cfg  config.PostgreSQL
}

func New(cfg config.PostgreSQL) *Storage {
	return &Storage{
		cfg: cfg,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	connectCtx, connectCtxCancel := context.WithTimeout(ctx, s.cfg.ConnectTimeout)
	defer connectCtxCancel()

	pool, err := pgxpool.New(connectCtx, s.cfg.URI)
	if err != nil {
		return errors.Newf("could not create postgresql connection pool: %w", err)
	}

	s.pool = pool

	return nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

type dbRecord interface {
	Scan(dest ...interface{}) error
}
