package transaction

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ExecTx(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}

	err = fn(tx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error rolling back tx: %w", err)
		}
		return fmt.Errorf("something went wrong with tx: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commiting tx: %w", err)
	}

	return nil
}
