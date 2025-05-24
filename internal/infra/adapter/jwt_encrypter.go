package adapter

import "github.com/golang-jwt/jwt/v4"

type JWTEncrypter struct {
	secret string
}

func NewJWTEncrypter(secret string) *JWTEncrypter {
	return &JWTEncrypter{
		secret: secret,
	}
}
func (j *JWTEncrypter) Encrypt(plainText string) (string, error) {
	claims := jwt.MapClaims{
		"data": plainText,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secret))

	if err != nil {
		return "", err
	}
	return signedToken, nil
}
