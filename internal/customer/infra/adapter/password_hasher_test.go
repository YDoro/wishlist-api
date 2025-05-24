package adapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ydoro/wishlist/internal/customer/infra/adapter"
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
