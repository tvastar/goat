package secrets_test

import (
	"context"
	"testing"

	"github.com/tvastar/goat"
	"github.com/tvastar/goat/secrets"
)

func TestVault(t *testing.T) {
	var v goat.EncrypterDecrypter = &secrets.Vault{
		EncryptPath: "transit/encrypt/goat",
		DecryptPath: "transit/decrypt/goat",
	}
	ctx := context.Background()
	data, err := v.Encrypt(ctx, []byte("hello"))
	if err != nil {
		t.Fatal("Failed to encrypt", err)
	}
	data, err = v.Decrypt(ctx, data)
	if err != nil || string(data) != "hello" {
		t.Fatal("Failed to decrypt", string(data), err)
	}
}

/*

sadly this does not work for some reason. It fails with the
following error:

     vault_test.go:17: could not enable secrets engine Error making API request.

        URL: POST https://127.0.0.1:54225/v1/sys/mounts/transit
        Code: 404. Raw Message:

        404 page not found

func createTestVault(t *testing.T) (*api.Config, string) {
	t.Helper()

	cluster := vault.NewTestCluster(t, nil, nil)
	t.Cleanup(cluster.Cleanup)
	cluster.Start()

	port := cluster.Cores[0].Listeners[0].Address.Port

	transport := cleanhttp.DefaultTransport()
	transport.TLSClientConfig = cluster.Cores[0].TLSConfig.Clone()
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}
	conf := api.DefaultConfig()
	conf.Address = fmt.Sprintf("https://127.0.0.1:%d", port)
	conf.HttpClient = &http.Client{Transport: transport}

	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatal("cannot create client", err)
	}
	client.SetToken(cluster.RootToken)
	err = client.Sys().Mount("transit", &api.MountInput{Type: "transit"})
	if err != nil {
		t.Fatal("could not enable secrets engine", err)
	}

	return conf, cluster.RootToken
}
*/
