package database

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

type Configuration struct {
	ConnectionURL      string
	MaxConnections     int
	MaxIdleConnections int
	MaxConnLifeTime    time.Duration
	MaxConnIdleTime    time.Duration
}

type Lifecycle interface {
	Close() error
}

type Client struct {
	Client *gorm.DB
}

var (
	db      *Client
	once    sync.Once
	initErr error
)

const (
	maxConnectionAttempts = 5
	retryDelay            = 2 * time.Second
)
