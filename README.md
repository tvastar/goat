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

## Goat integration

Goat requires a session storage, a token storage and an encryption
engine.

A Redis-based session storage engine is implemented in the ./sessions
package.

A [entgo.io](https://entgo.io/) based token storage is implemented in
the ./tokens package.

A [vault](https://github.com/hashicorp/vault/) based encryption engine
is implemented in the ./secrets package (where vault is used as a
secrets transit engine rather than actual storage).

All of these are optional with custom implementations possible. An
[example](https://github.com/tvastar/goat/blob/master/cmd/goat/goat.go)
which pulls these together is provided along with its [assocciated config](https://github.com/tvastar/goat/blob/master/cmd/goat/config.yaml.sample)

This example requires a local redis and a local vault. The token
storage uses an in-memory sqllite3 DB.  The example and the associated
redis/vault containers can be launched via the `./scripts/run.sh`
script.

For the example to work, the
[config.yaml.sample](https://github.com/tvastar/goat/blob/master/cmd/goat/config.yaml.sampl)
file must be updated to specify a valid provider and then one should
visit the following url in the browser
[http://localhost:8085/sheets/url?redirect_url=http://www.google.com](http://localhost:8085/sheets/url?redirect_url=http://www.google.com).

The URL above should walk through the consent flow and end up on
google.  At this point, the current token can be fetched by visiting
[this url](http://localhost:8085/sheets/token)

## Sample Config file

The following is a sample config file for this service:

```yaml
httpport: 8085
providers:
  - name: sheets
    paths:
      consent: /sheets/url
      code: /sheets/code
      setrefreshtoken: /sheets/setRefreshToken
      getaccesstoken: /sheets/token
    config:
      endpoint:
        authurl: https://accounts.google.com/o/oauth2/auth
        tokenurl: https://oauth2.googleapis.com/token
      clientid: <your google API project client_id>
      clientsecret: <your google API project client_secret>
      redirecturl: http://localhost:8085/sheets/code
      scopes: ["email"]
tokens:
  dbsource: "file:ent?mode=memory&cache=shared&_fk=1"
  dbtype: sqlite3
sessions:
  ttl: 30s
  options:
    addr: "localhost:6379"
secrets:
  encryptpath: transit/encrypt/goat
  decryptpath: transit/decrypt/goat
```

## Running the tests

Tests can be run locally via the `./scripts/test.sh` script (which
launches the required docker images and such.

The following manual process can be used to setup the vault container:

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

