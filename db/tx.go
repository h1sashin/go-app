package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TxKey string

const txKey TxKey = "tx"

func InjectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func ExtractTx(ctx context.Context) (pgx.Tx, bool) {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx, true
	}
	return nil, false
}
