// Package secrets implements goat secret encryption/decryption.
package secrets

import (
	"context"
	"encoding/base64"
	"sync"
	vault "github.com/hashicorp/vault/api"
)

// Vault uses the vault transit engine to encrypt/decrypt tokens.
//
// See https://learn.hashicorp.com/vault/encryption-as-a-service/eaas-transit
//
// Vault config is optional. By default, vault can be configured via
// the environment variables VAULT_ADDR and VAULT_TOKEN.  But explicit
// configuration is possible by providing a Config with Address
// specified as well as by providing a Token.
//
// The EncryptPath and DecryptPath must match the secrets engine
// path.  Typically, vault is configured via:
//
//     vault secrets enable transit
//     vault write -f transit/keys/goat
//
// If so, the EncryptPath and DecryptPath would be
// transit/encrypt/goat and transit/decrypt/goat respectively.
//
// The transit engine can be set up on a different path:
//
//      vault secret enable -path=foo transit
//
// Similarly, the keyring can also be something other than "goat":
//
//      vault write -f foo/keys/my-key
//
// In this case, the EncryptPath and DecryptPath would become
// foo/encrypt/my-key and foo/decrypt/my-key respectively.
type Vault struct {
	*vault.Config
	Token string
	EncryptPath string
	DecryptPath string
	c *vault.Client
	mu sync.Mutex
}

// Encrypt encrypts using vault.
func (v *Vault) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	client, err := v.initClient()
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{"plaintext": data}
	secret, err := client.Logical().Write(v.EncryptPath, body)
	if err != nil {
		return nil, err
	}
	return []byte(secret.Data["ciphertext"].(string)), nil
}

// Decrypt decrypts using vault.
func (v *Vault) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	client, err := v.initClient()
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{"ciphertext": string(data)}
	secret, err := client.Logical().Write(v.DecryptPath, body)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(secret.Data["plaintext"].(string))
}

func (v *Vault) initClient() (*vault.Client, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.c != nil {
		return v.c, nil
	}
	c, err := vault.NewClient(v.Config)
	if err == nil {
		v.c = c
		if v.Token != "" {
			c.SetToken(v.Token)
		}
	}
	return c, err
}






