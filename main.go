package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wreulicke/go-advisory-lock-sandbox/internal/db"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", "root", "root", "localhost", 15432, "postgres")
	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return pool, nil
}

func main() {
	p, _ := NewPool(context.Background())
	queries := db.New(p)

	tx, err := p.Begin(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	ok, err := queries.WithTx(tx).TryAdvisoryLock(context.TODO(), "test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ok)
	if !ok {
		return
	}
	// do long-term tasks with lock
	time.Sleep(5 * time.Second)

	// unlock by commit
	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("unlock")

	time.Sleep(5 * time.Second)
}
