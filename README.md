# Goat

Goat is a general purpose oauth2 access token server.  It provides a
way for servers and clients to securely maintain Oauth refresh
tokens.

The typical flow is as follows:

1. A browser client makes an authenticated HTTP request to Goat,
requesting credentials for a specific Oauth provider (such as
"google").   This provider must be configured within Goat with the
Oauth credentials.  In addition this oauth provider must be configured
with redirect URLs to Goat.

2. Goat checks if it has the oauth credentials for the current
authenticated user.  If it does, it responds with an oauth Token JSON
(but without the refresh_token part).  Note: if the token stored with
Goat is expired, it refreshes it automatically.

3. Goat returns an empty hash if it doensn't have the credentials.  In
this case, the browser client needs to initiate an Oauth "consent"
flow with the oauth provider.  The browser client redirects the user
to the Goat server consent endpoint for that provider.

4. Goat calculates a nonce to prevent CSRF attacks, saves the redirect
URL to the session store and then redirects the browser to the actual
consent page of the oauth provider.

5. When the user completes the consent flow of the oauuth provider,
the oauth provider redirects to the configured redirect URL, which
should be the "auth code" endpoint on the Goat server.

6. The Goat server handles the auth code provided by the oauth
provider and exchanges it to get the token.  It encrypts this token
using an encryption engine and then saves it within the token
storage.  It then redirects the browser back to the redirect URL
configured by the client in step 3.

7. The client detects that the consent flow has completed
successfully and then initiates step 1 again, which will now succeed.

## Goat configuration

Goat requires a session storage, a token storage and an encryption
engine.

## Running the tests

The ./secrets folder requires vault to be running:

```bash
docker run --rm --cap-add=IPC_LOCK -e VAULT_DEV_ROOT_TOKEN_ID=hello --name=dev-vault -p 8200:8200 vault:1.5.0
```

In addition, the example expects vault to be configured as a transit
engine and a keyring:

```bash
docker exec -it $(docker ps -q -f name=dev-vault) sh -c 'VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=hello vault secrets
enable transit'
docker exec -it $(docker ps -q -f name=dev-vault) sh -c 'VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=hello vault write -f transit/keys/goat'
```


Now the tests can be run using:

```bash
VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=hello go test ./...
```

## Example

See [this example
server](https://github.com/tvastar/goat/blob/master/cmd/goat/goat.go)
for how to set up your service.

