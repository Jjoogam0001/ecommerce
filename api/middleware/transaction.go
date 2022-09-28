package middleware

import (
	"context"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

// ErrNoTransactionContext defines the error when echo context is not a transaction context.
var ErrNoTransactionContext = errors.New("current context is not a transaction context, check using of the middleware")

// TxProvider is the definition of a Tx factory. This contract links the transaction begin behavior of pgxpool.Pool and pgx.Tx.
type TxProvider interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type transactionCtx struct {
	echo.Context
	tx pgx.Tx
}

// FromTransactionContext returns the transaction from the context if exists.
func FromTransactionContext(ctx echo.Context) (pgx.Tx, error) {
	t, ok := ctx.(*transactionCtx)
	if !ok {
		return nil, ErrNoTransactionContext
	}
	return t.tx, nil
}

// Transaction middleware to handle per-request database transactions.
func Transaction(txProvider TxProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			tx, err := txProvider.Begin(ctx)
			if err != nil {
				return err
			}

			if err := next(&transactionCtx{c, tx}); err != nil {
				defer func() {
					if err := tx.Rollback(ctx); err != nil {
						c.Logger().Error(err)
					}
				}()

				return err
			}

			return tx.Commit(ctx)
		}
	}
}
