#!/usr/bin/env bash

set -ex

docker run --rm -d --cap-add=IPC_LOCK -e VAULT_DEV_ROOT_TOKEN_ID=hello --name=dev-vault -p 8200:8200 vault:1.5.0 
trap 'docker stop $(docker ps -q -f name=dev-vault)' EXIT INT

sleep 3

docker exec -it $(docker ps -q -f name=dev-vault) sh -c 'VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=hello vault secrets enable transit'

docker exec -it $(docker ps -q -f name=dev-vault) sh -c 'VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=hello vault write -f transit/keys/goat'

VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=hello go test ./...
