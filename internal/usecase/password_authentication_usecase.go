package usecase

import (
	"context"
	"encoding/json"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
	"github.com/ydoro/wishlist/internal/presentation/inputs"
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

func (p *PasswordAuthenticationUseCase) Authenticate(ctx context.Context, credentials any) (string, error) {
	pwdAuth, ok := credentials.(inputs.PwdAuth)
	if !ok {
		return "", &e.ValidationError{
			Field: "credentials",
			Err:   "Invalid credentials type",
		}
	}
	// retrieve the user by email
	user, err := p.UserGetter.GetByEmail(ctx, pwdAuth.Email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", e.NewAuthenticationError(domain.AuthMethodPassword)
	}

	// compare the password with the hash
	err = p.HashComparer.Compare(user.Password, pwdAuth.Password)
	if err != nil {
		return "", e.NewAuthenticationError(domain.AuthMethodPassword)
	}
	// if the password matches, sign a JWT token and return it
	outgoing := &domain.OutgoingCustomer{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	data, _ := json.Marshal(outgoing)
	token, err := p.Encrypter.Encrypt(string(data))

	if err != nil {
		return "", err
	}

	return token, nil
}
