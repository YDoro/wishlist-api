//go:generate mockgen --build_flags=--mod=mod -destination=../../mock/domain/auth_mock.go -package=mocks . Authenticator,Authorizer

package domain

import "context"

type AuthMethod string

const (
	AuthMethodPassword AuthMethod = "password"
)

type Authenticator interface {
	Authenticate(ctx context.Context, credentials any) (string, error)
}

type Authorizer interface {
	Authorize(ctx context.Context, token string) (string, error)
}
