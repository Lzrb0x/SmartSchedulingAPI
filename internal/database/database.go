package database

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Lzrb0x/SmartSchedulingAPI/internal/config"
)

func Connect(ctx context.Context, cfg config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", cfg.URL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleSec) * time.Second)

	return db, nil
}
