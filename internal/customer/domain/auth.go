package domain

import "context"

type AuthMethod string

const (
	AuthMethodPassword AuthMethod = "password"
)

type Authenticator interface {
	Authenticate(ctx context.Context, credentials any) (string, error)
}
