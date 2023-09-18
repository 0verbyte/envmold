package mold

import (
	"strings"
	"testing"
)

func TestSecretManagerMock(t *testing.T) {
	var moldTemplate = `
- name: foo
  value: 'mock("test/foo")'
  type: string
  required: false

- name: debug
  value: true
  type: boolean
  required: false
`

	m, err := New(strings.NewReader(moldTemplate), nil)
	if err != nil {
		t.Errorf("Failed to create new mold: %v", err)
		return
	}

	envVarKey := "foo"
	v, err := m.GetVariable(envVarKey)
	if err != nil {
		t.Errorf("Failed to get environment variable from mold: %s", envVarKey)
		return
	}

	secretManager, err := checkAndUseSecretManager(v.Value.(string))
	secretManagerMock := secretManager.(*SecretManagerMock)
	secretManagerMock.Seed()

	if err != nil {
		t.Errorf("Got error: %v", err)
		return
	}

	got, err := secretManager.GetValue(v.Value.(string))
	if err != nil {
		t.Errorf("Got error trying to read value from secret manager: %v", err)
		return
	}

	if got != "mock_bar" {
		t.Errorf("Expected mock_bar from secret manager mock, got: %s", got)
		return
	}
}
