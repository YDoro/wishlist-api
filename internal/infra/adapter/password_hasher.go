package adapter

import "golang.org/x/crypto/bcrypt"

type PasswordHasher struct {
	salt int
}

func NewPasswordHasher(salt int) *PasswordHasher {
	return &PasswordHasher{
		salt: salt,
	}
}

// Hash generates a hashed version of the given password using the bcrypt algorithm.
// The salt value, which determines the computational cost of the hashing process, is
// provided by the PasswordHasher instance. It returns the hashed password as a string
// and any potential error encountered during the hashing process.
func (h *PasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.salt)
	return string(hash), err
}

func (h *PasswordHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
