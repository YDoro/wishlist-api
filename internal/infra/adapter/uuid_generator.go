package adapter

import "github.com/google/uuid"

type UUIDGenerator struct {
}

func (UUIDGenerator) Generate() (string, error) {
	uid, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	return uid.String(), nil
}
