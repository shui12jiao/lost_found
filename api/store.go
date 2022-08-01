package api

import (
	"context"
	"database/sql"
	"fmt"
	"lost_found/db/sqlc"
)

//提供执行查询和事务的接口
type Store interface {
	sqlc.Querier
}

//Store的实现
type SQLStore struct {
	*sqlc.Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: sqlc.New(db),
	}
}

//执行fn事务
func (store *SQLStore) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
