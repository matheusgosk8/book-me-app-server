package repository

import (
	"context"

	"github.com/matheusgosk8/book-me-server/ent"
)

// RunInTx executes the provided function inside an ent transaction.
// The function receives the transaction `tx` to run repository operations.
// On error the transaction is rolled back, otherwise committed.
func RunInTx(ctx context.Context, client *ent.Client, fn func(context.Context, *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
