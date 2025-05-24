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
