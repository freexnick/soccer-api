package repository

import "context"

type Repositories struct {
	User     User
	Team     Team
	Player   Player
	Transfer Transfer
}

type RepoFunc func(repos Repositories) error

type TransactionManager interface {
	Execute(ctx context.Context, fn RepoFunc) error
}
