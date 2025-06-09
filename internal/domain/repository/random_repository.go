package repository

import "context"

type RandomData interface {
	FirstName(ctx context.Context) string
	LastName(ctx context.Context) string
	Age(ctx context.Context, min, max int) int
}
