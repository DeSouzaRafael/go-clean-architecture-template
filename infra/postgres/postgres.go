package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
	defaultPingTimeout  = 10 * time.Second
)

type Options struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(opts Options) (*Postgres, error) {
	gormDB, err := gorm.Open(postgres.Open(opts.URL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if opts.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(opts.MaxOpenConns)
	}
	if opts.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(opts.MaxIdleConns)
	}
	if opts.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(opts.ConnMaxLifetime)
	}

	for i := 0; i < defaultConnAttempts; i++ {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), defaultPingTimeout)
		err = sqlDB.PingContext(pingCtx)
		pingCancel()
		if err == nil {
			break
		}
		if i < defaultConnAttempts-1 {
			time.Sleep(defaultConnTimeout)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect after %d attempts: %w", defaultConnAttempts, err)
	}

	return &Postgres{DB: gormDB}, nil
}

func (p *Postgres) Close() {
	if p.DB != nil {
		if sqlDB, err := p.DB.DB(); err == nil {
			sqlDB.Close()
		}
	}
}
