package service

import (
	"context"
	"lost_found/db/sqlc"
)

type Transcation interface {
	CompleteFoundTx(ctx context.Context, param CompleteFoundTxParams) (sqlc.Match, error)
	CompleteLostTx(ctx context.Context, param CompleteFoundTxParams) (sqlc.Match, error)
}

//提交拾取物已寻回
type CompleteFoundTxParams struct {
	AddMatchParam sqlc.AddMatchParams `json:"addMatchParams"`
	ID            int32               `json:"id"`
}

func (store *SQLStore) CompleteFoundTx(ctx context.Context, param CompleteFoundTxParams) (sqlc.Match, error) {
	var match sqlc.Match
	err := store.execTx(ctx, func(q *sqlc.Queries) error {
		var err error
		match, err = q.AddMatch(ctx, param.AddMatchParam)
		if err != nil {
			return err
		}
		err = q.DeleteFound(ctx, param.ID)
		return err
	})

	return match, err
}

//提交遗失物已寻回
type CompleteLostTxParams = CompleteFoundTxParams

func (store *SQLStore) CompleteLostTx(ctx context.Context, param CompleteFoundTxParams) (sqlc.Match, error) {
	var match sqlc.Match
	err := store.execTx(ctx, func(q *sqlc.Queries) error {
		var err error
		match, err = q.AddMatch(ctx, param.AddMatchParam)
		if err != nil {
			return err
		}
		err = q.DeleteLost(ctx, param.ID)
		return err
	})

	return match, err
}
