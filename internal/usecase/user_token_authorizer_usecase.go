package usecase

import (
	"context"

	"github.com/ydoro/wishlist/internal/domain"
)

type UserTokenAuthorizerUseCase struct {
	TokenDecoder domain.Decrypter
}

func NewUserTokenAuthorizerUseCase(decrypter domain.Decrypter) *UserTokenAuthorizerUseCase {
	return &UserTokenAuthorizerUseCase{
		TokenDecoder: decrypter,
	}
}

func (u *UserTokenAuthorizerUseCase) Authorize(ctx context.Context, token string) (string, error) {
	return u.TokenDecoder.Decrypt(token)
}
