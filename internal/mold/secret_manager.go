package mold

import (
	"errors"
	"strings"
)

var (
	ErrSecretManagerKeyDoesNotExist = errors.New("key does not exist")
	ErrSecretManagerNotFound        = errors.New("value does not contain a secret manager")
)

type SecretManager interface {
	GetValue(string) (string, error)
}

var secretManagerKeysToTypes = map[string]SecretManager{
	"mock": &SecretManagerMock{},
}

func checkAndUseSecretManager(s string) (SecretManager, error) {
	for secretManagerKey, secretManager := range secretManagerKeysToTypes {
		if strings.HasPrefix(s, secretManagerKey) {
			return secretManager, nil
		}
	}
	return nil, ErrSecretManagerKeyDoesNotExist
}
