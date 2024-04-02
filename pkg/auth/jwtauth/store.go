package jwtauth

import (
	"context"
	"time"
)

type Storer interface {
	Set(ctx context.Context, tokenString string, expireation time.Duration) error

	Check(ctx context.Context, tokenString string) (bool, error)

	Close() error
}
