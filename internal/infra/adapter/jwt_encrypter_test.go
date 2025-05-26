package adapter_test

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/ydoro/wishlist/internal/infra/adapter"
)

func TestJWTEncrypter_Encrypt(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		plainText string
		wantErr   bool
	}{
		{
			name:      "successful encryption",
			secret:    "mysecret",
			plainText: "hello world",
			wantErr:   false,
		},
		{
			name:      "empty plain text",
			secret:    "mysecret",
			plainText: "",
			wantErr:   false,
		},
		{
			name:      "empty secret",
			secret:    "",
			plainText: "hello world",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypter := adapter.NewJWTEncrypter(tt.secret)
			token, err := encrypter.Encrypt(tt.plainText)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.NotEqual(t, tt.plainText, token)

			// TODO - Verify using Decryption method
			// Verify the token can be parsed and contains correct data
			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(tt.secret), nil
			})

			assert.NoError(t, err)
			assert.True(t, parsedToken.Valid)

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			assert.True(t, ok)
			assert.Equal(t, tt.plainText, claims["data"])
		})
	}
}

func TestJWTEncrypter_Decrypt(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		plainText string
		wantErr   bool
	}{
		{
			name:      "successful decryption",
			secret:    "mysecret",
			plainText: "hello world",
			wantErr:   false,
		},
		{
			name:      "empty plain text",
			secret:    "mysecret",
			plainText: "",
			wantErr:   false,
		},
		{
			name:      "wrong secret",
			secret:    "wrongsecret",
			plainText: "hello world",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypter := adapter.NewJWTEncrypter("mysecret")
			token, err := encrypter.Encrypt(tt.plainText)
			assert.NoError(t, err)

			decrypter := adapter.NewJWTEncrypter(tt.secret)
			decrypted, err := decrypter.Decrypt(token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, decrypted)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.plainText, decrypted)
		})
	}
}
