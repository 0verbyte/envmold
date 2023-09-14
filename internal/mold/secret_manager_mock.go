package mold

import (
	"strings"
)

const (
	SecretManagerMockIdentifier = "mock"
)

type SecretManagerMock struct {
	values map[string]string
}

func (s *SecretManagerMock) Seed() {
	s.values = map[string]string{}
	s.values["test/foo"] = "mock_bar"
}

// Mold passes the value as defined in the template, including the secret manager identifier.
// Example: mock("test/foo_bar").
func getKeyFromMoldValue(s string) string {
	s = strings.TrimPrefix(s, SecretManagerMockIdentifier)
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		s = s[1 : len(s)-1]
	}
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		s = s[1 : len(s)-1]
	}
	return s
}

func (s *SecretManagerMock) GetValue(key string) (string, error) {
	if v, ok := s.values[getKeyFromMoldValue(key)]; ok {
		return v, nil
	}
	return "", ErrSecretManagerKeyDoesNotExist
}
