package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Client struct {
	*sqlx.DB
	queryTimeOut time.Duration
}

func New(cfg Config) (*Client, error) {
	d, err := sqlx.Open(cfg.Driver, cfg.URL)
	if err != nil {
		return nil, err
	}

	if err = d.Ping(); err != nil {
		return nil, err
	}

	d.SetMaxIdleConns(cfg.MaxIdleConns)
	d.SetMaxOpenConns(cfg.MaxOpenConns)
	d.SetConnMaxLifetime(cfg.ConnMaxLifeTime)

	return &Client{DB: d, queryTimeOut: cfg.MaxQueryTimeoutInMS}, err
}

func (c Client) WithTimeout(ctx context.Context, op func(ctx context.Context) error) (err error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.queryTimeOut)
	defer cancel()

	return op(ctxWithTimeout)
}

// Handling transactions: https://stackoverflow.com/a/23502629/8244298
func (c Client) WithTxn(ctx context.Context, txnOptions sql.TxOptions, txFunc func(*sqlx.Tx) error) (err error) {
	txn, err := c.BeginTxx(ctx, &txnOptions)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = errors.Errorf("%s", p)
			}
			err = txn.Rollback()
			panic(p)
		} else if err != nil {
			if rlbErr := txn.Rollback(); err != nil {
				err = fmt.Errorf("rollback error: %s while executing: %w", rlbErr, err)
			} else {
				err = fmt.Errorf("rollback: %w", err)
			}
			err = fmt.Errorf("rollback: %w", err)
		} else {
			err = txn.Commit()
		}
	}()

	err = txFunc(txn)
	return err
}
