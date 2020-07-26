// This command demonstrates how to run a goat server.
//
// Goat requires a tokens storage, a session storage and an
// enccryption engine.  This example uses the ent-based token storage
// (against sqlite3), a REDIS based sessions store and a vault based
// secret encryption engine.
//
// All of this can be configured neatly via the sample yaml in this
// directory.
//
// Launch this program with the sample yaml as an argument and then
// hit the following URL in the browser
// http://localhost:8085/sheets/url?redirect_url=http://www.google.com
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tvastar/goat"
	"github.com/tvastar/goat/secrets"
	"github.com/tvastar/goat/sessions"
	"github.com/tvastar/goat/tokens"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers []*goat.Provider
	HTTPPort  int
	Tokens    tokens.Tokens
	Secrets   secrets.Vault
	Sessions  sessions.Redis
}

func main() {
	var c Config
	if len(os.Args) != 2 {
		log.Fatal("Usage: goat yaml_config_file")
	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("reading config", err)
	}
	if err = yaml.Unmarshal(data, &c); err != nil {
		log.Fatal("parsing yaml", err)
	}

	s := &http.Server{
		Addr: fmt.Sprintf(":%d", c.HTTPPort),
		Handler: goat.Handler{
			Providers:         c.Providers,
			AuthenticatedUser: authenticatedUser,
			Tokens:            &c.Tokens,
			Sessions:          &c.Sessions,
			Secrets:           &c.Secrets,
		},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func authenticatedUser(r *http.Request) (string, error) {
	return "foo", nil
}
