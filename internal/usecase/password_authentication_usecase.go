package usecase

import (
	"context"
	"encoding/json"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/pkg/presentation/errors"
	"github.com/ydoro/wishlist/pkg/presentation/inputs"
)

type PasswordAuthenticationUseCase struct {
	HashComparer domain.HashComparer
	UserGetter   domain.GetCustomerByEmailRepository
	Encrypter    domain.Encrypter
}

func NewPasswordAuthenticationUseCase(comparer domain.HashComparer, userGetter domain.GetCustomerByEmailRepository, encrypter domain.Encrypter) *PasswordAuthenticationUseCase {
	return &PasswordAuthenticationUseCase{
		HashComparer: comparer,
		UserGetter:   userGetter,
		Encrypter:    encrypter,
	}
}

func (p *PasswordAuthenticationUseCase) Authenticate(ctx context.Context, credentials inputs.PwdAuth) (string, error) {
	// retrieve the user by email
	user, err := p.UserGetter.GetByEmail(ctx, credentials.Email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", e.NewAuthenticationError(domain.AuthMethodPassword)
	}

	// compare the password with the hash
	err = p.HashComparer.Compare(user.Password, credentials.Password)
	if err != nil {
		return "", e.NewAuthenticationError(domain.AuthMethodPassword)
	}
	// if the password matches, sign a JWT token and return it
	data, _ := json.Marshal(user)
	token, err := p.Encrypter.Encrypt(string(data))

	if err != nil {
		return "", err
	}

	return token, nil
}
