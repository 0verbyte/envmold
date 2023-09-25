package mold

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault-client-go"
)

type SecretManagerVault struct {
}

var vaultClient *vault.Client

var (
	ErrVaultKeyMissing = "%s is missing"
)

func initVaultClient() error {
	client, err := vault.New(
		vault.WithAddress("127.0.0.1:8200"),
		vault.WithRequestTimeout(5*time.Second),
	)
	if err != nil {
		return err
	}

	vaultClient = client

	return nil
}

func (s *SecretManagerVault) GetValue(key string) (string, error) {
	if vaultClient == nil {
		if err := initVaultClient(); err != nil {
			return "", err
		}
	}
	ctx := context.Background()
	secret, err := vaultClient.Secrets.KvV2Read(ctx, key)
	if err != nil {
		return "", err
	}

	if v, ok := secret.Data.Data[key].(string); ok {
		return v, nil
	}

	return "", fmt.Errorf(ErrVaultKeyMissing, key)
}
