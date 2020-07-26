package tokens

import (
	"context"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/tvastar/goat/tokens/ent"
	"github.com/tvastar/goat/tokens/ent/tokens"
)

type Tokens struct {
	DBSource string
	DBType   string
	c        *ent.Client
	mu       sync.Mutex
}

func (t *Tokens) Get(ctx context.Context, provider, user string) ([]byte, error) {
	client, err := t.getClient(ctx)
	if err != nil {
		return nil, err
	}

	tok, err := client.Tokens.Query().
		Where(tokens.Provider(provider), tokens.User(user)).
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, os.ErrNotExist
	}

	if err != nil {
		return nil, err
	}
	return []byte(tok.Token), nil
}

func (t *Tokens) Set(ctx context.Context, provider, user string, token []byte) error {
	client, err := t.getClient(ctx)
	if err != nil {
		return err
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	count, err := client.Tokens.Update().
		Where(tokens.Provider(provider), tokens.User(user)).
		SetToken(string(token)).
		Save(ctx)
	if err == nil && count == 0 {
		_, err = client.Tokens.Create().
			SetProvider(provider).
			SetUser(user).
			SetToken(string(token)).
			Save(ctx)
	}

	if err == nil {
		return tx.Commit()
	}

	if rerr := tx.Rollback(); rerr != nil {
		return errors.Wrap(err, rerr.Error())
	}

	return err
}

func (t *Tokens) getClient(ctx context.Context) (*ent.Client, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.c != nil {
		return t.c, nil
	}

	client, err := ent.Open(t.DBType, t.DBSource)
	if err == nil {
		err = client.Schema.Create(ctx)
		if err == nil {
			t.c = client
		} else {
			client.Close()
		}
	}
	return t.c, err
}

func (t *Tokens) Close() error {
	if t.c != nil {
		return t.c.Close()
	}
	return nil
}
