package postgresql

import (
	"github.com/ZdravkoGyurov/myx-image-api/pkg/errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

const (
	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"
)

func dbError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.Newf("%s: %w", err, errors.ErrEntityNotFound)
	}

	var pgxErr *pgconn.PgError
	if ok := errors.As(err, &pgxErr); ok {
		switch pgxErr.Code {
		case uniqueViolation:
			return errors.Newf("%s: %w", err, errors.ErrEntityAlreadyExists)
		case foreignKeyViolation:
			return errors.Newf("%s: %w", err, errors.ErrRefEntityViolation)
		}
	}

	return err
}
