//go:generate mockgen --build_flags=--mod=mod -destination=../../mock/domain/criptography_mock.go -package=mocks . Hasher,HashComparer,Encrypter,Decrypter

package domain

type Hasher interface {
	Hash(password string) (string, error)
}

type HashComparer interface {
	Compare(hashedPassword, password string) error
}

type Encrypter interface {
	Encrypt(plainText string) (string, error)
}

type Decrypter interface {
	Decrypt(cipherText string) (string, error)
}
