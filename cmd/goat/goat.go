package main

import (
	"context"
	"fmt"
	"github.com/tvastar/goat"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Providers []*goat.Provider
	HTTPPort  int
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
			Tokens:            &memstore{items: map[[2]string][]byte{}},
			Sessions:          &memstore{items: map[[2]string][]byte{}},
			Secrets:           &secrets{},
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

type memstore struct {
	items map[[2]string][]byte
}

func (m *memstore) Get(_ context.Context, provider, key string) ([]byte, error) {
	if v, ok := m.items[[2]string{provider, key}]; ok {
		return v, nil
	}
	return nil, os.ErrNotExist
}

func (m *memstore) Set(_ context.Context, provider, key string, data []byte) error {
	m.items[[2]string{provider, key}] = data
	return nil
}

type secrets struct{}

func (s *secrets) Encrypt(_ context.Context, data []byte) ([]byte, error) {
	return data, nil
}

func (s *secrets) Decrypt(_ context.Context, data []byte) ([]byte, error) {
	return data, nil
}
