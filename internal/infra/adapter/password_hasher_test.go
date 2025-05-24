package adapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ydoro/wishlist/internal/infra/adapter"
)

func TestPasswordHasher_Hash(t *testing.T) {
	cases := []struct {
		name     string
		password string
		salt     int64
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "mysecretpassword",
			salt:     10,
			wantErr:  false,
		},
		{
			name:     "should not return error on empty password",
			password: "",
			salt:     10,
			wantErr:  false,
		},
		{
			name:     "invalid salt",
			password: "mysecretpassword",
			salt:     int64(999999999999999999), // This is an invalid value for bcrypt
			wantErr:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			hasher := adapter.NewPasswordHasher(int(c.salt))
			hash, err := hasher.Hash(c.password)

			if c.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, c.password, hash)
			}

		})
	}

}

func TestPasswordHasher_Compare(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		setupPassword string
		wantErr       bool
	}{
		{
			name:          "matching passwords",
			setupPassword: "correctpassword",
			password:      "correctpassword",
			wantErr:       false,
		},
		{
			name:          "non-matching passwords",
			setupPassword: "correctpassword",
			password:      "wrongpassword",
			wantErr:       true,
		},
		{
			name:          "empty password should not match hash",
			setupPassword: "somepassword",
			password:      "",
			wantErr:       true,
		},
		{
			name:          "empty setup password",
			setupPassword: "",
			password:      "somepassword",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hasher := adapter.NewPasswordHasher(10)

			hashedPassword, err := hasher.Hash(tt.setupPassword)
			assert.NoError(t, err)

			err = hasher.Compare(hashedPassword, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
