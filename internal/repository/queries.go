package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (q *Queries) ResetWorkflowRun(ctx context.Context, pool *pgxpool.Pool, id uuid.UUID) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := q.WithTx(tx)

	if err = qtx.UpdateWorkflowRunStatusNull(ctx, id); err != nil {
		return err
	}
	if err = qtx.UpdateJobRunStatusNull(ctx, id); err != nil {
		return err
	}
	if err = qtx.UpdateStepRunStatusNull(ctx, id); err != nil {
		return err
	}
	if err = qtx.DeleteCommandOutputByWorkflowId(ctx, id); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
