//go:generate mockgen --build_flags=--mod=mod -destination=../../mock/domain/cacge_mock.go -package=mocks -source ./cache.go

package domain

import (
	"time"

	"golang.org/x/net/context"
)

type Cache interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, keys ...string) error
}
