package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultMaxPoolSize  = 10
	defaultConnAttempts = 10
	defaultConnTimeout  = 10 * time.Second
)

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(url string) (*Postgres, error) {

	gormDB, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to the database using GORM: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(defaultMaxPoolSize)
	sqlDB.SetMaxIdleConns(defaultMaxPoolSize / 2)
	sqlDB.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), defaultConnTimeout)
	defer cancel()

	for i := 0; i < defaultConnAttempts; i++ {
		if err = sqlDB.PingContext(ctx); err == nil {
			break
		}
		if i < defaultConnAttempts-1 {
			time.Sleep(defaultConnTimeout)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database after %d attempts: %w", defaultConnAttempts, err)
	}

	return &Postgres{DB: gormDB}, nil
}

func (p *Postgres) Close() {
	if p.DB != nil {
		sqlDB, err := p.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}
