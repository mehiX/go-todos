package db

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Conn(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

type dbConnFn func(ctx context.Context, dsn string) (*sql.DB, error)

func ConnWithRetry(f dbConnFn, retries int, base, cap time.Duration) dbConnFn {
	backoff := base

	return func(ctx context.Context, dsn string) (*sql.DB, error) {
		for r := 0; ; r++ {
			conn, err := f(ctx, dsn)
			if r >= retries || err == nil {
				return conn, err
			}

			if backoff > cap {
				backoff = cap
			}

			jitter := rand.Int63n(int64(base))
			wait := backoff + time.Duration(jitter)

			fmt.Printf("DB connection number %d failed. Waiting %v before retrying\n", r+1, wait)

			select {
			case <-time.After(wait):
				backoff <<= 1
			case <-ctx.Done():
				return nil, ctx.Err()
			}

		}
	}
}
