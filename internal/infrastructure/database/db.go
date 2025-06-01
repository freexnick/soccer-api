package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func initDB(ctx context.Context, conf Configuration) (*Client, error) {
	silentGormLogger := gormlogger.Default.LogMode(gormlogger.Silent)

	var gormAPIDB *gorm.DB
	var lastAttemptErr error

	for attempt := 1; attempt <= maxConnectionAttempts; attempt++ {
		tempDB, err := gorm.Open(postgres.Open(conf.ConnectionURL), &gorm.Config{
			Logger: silentGormLogger,
		})

		if err != nil {
			lastAttemptErr = fmt.Errorf("gorm.Open (attempt %d): %w", attempt, err)
			if attempt < maxConnectionAttempts {
				time.Sleep(retryDelay)
				continue
			}
			break
		}

		sqlDB, err := tempDB.DB()
		if err != nil {
			lastAttemptErr = fmt.Errorf("getting underlying sql.DB (attempt %d): %w", attempt, err)
			if attempt < maxConnectionAttempts {
				time.Sleep(retryDelay)
				continue
			}
			break
		}

		sqlDB.SetMaxIdleConns(conf.MaxIdleConnections)
		sqlDB.SetMaxOpenConns(conf.MaxConnections)
		sqlDB.SetConnMaxLifetime(conf.MaxConnLifeTime)
		sqlDB.SetConnMaxIdleTime(conf.MaxConnIdleTime)

		pingErr := sqlDB.PingContext(ctx)
		if pingErr == nil {
			gormAPIDB = tempDB
			return &Client{Client: gormAPIDB}, nil
		}

		lastAttemptErr = fmt.Errorf("pinging database (attempt %d): %w", attempt, pingErr)
		if attempt < maxConnectionAttempts {
			time.Sleep(retryDelay)
		}
	}

	if lastAttemptErr != nil {
		return nil, fmt.Errorf("failed to connect to GORM database after %d attempts: %w", maxConnectionAttempts, lastAttemptErr)
	}

	return nil, errors.New("failed to connect to GORM database after all attempts (unexpected state)")
}

func New(ctx context.Context, conf Configuration) (*Client, error) {
	once.Do(func() {
		db, initErr = initDB(ctx, conf)
	})

	return db, initErr
}

func (c *Client) Close() error {
	if c == nil || c.Client == nil {
		return nil
	}

	sqlDB, err := c.Client.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB for GORM client for closing: %w", err)
	}

	return sqlDB.Close()
}
