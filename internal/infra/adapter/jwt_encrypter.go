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

func (j *JWTEncrypter) Decrypt(cipherText string) (string, error) {
	token, err := jwt.Parse(cipherText, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["data"].(string), nil
	}

	return "", jwt.NewValidationError("invalid token", jwt.ValidationErrorSignatureInvalid)
}
