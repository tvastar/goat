package tokens_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/tvastar/goat"
	"github.com/tvastar/goat/tokens"

	_ "github.com/mattn/go-sqlite3"
)

func TestTokens(t *testing.T) {
	var tok goat.TokenStore = &tokens.Tokens{
		DBSource: "file:ent?mode=memory&cache=shared&_fk=1",
		DBType:   "sqlite3",
	}
	defer closeIt(tok)

	ctx := context.Background()
	data, err := tok.Get(ctx, "provider", "user")
	if !os.IsNotExist(err) {
		t.Fatal("Get() without Set()", err, string(data))
	}

	err = tok.Set(ctx, "provider", "user", []byte("foo1"))
	if err != nil {
		t.Fatal("First Set()", err)
	}

	err = tok.Set(ctx, "provider", "user", []byte("foo2"))
	if err != nil {
		t.Fatal("Second Set()", err)
	}

	data, err = tok.Get(ctx, "provider", "user")
	if err != nil || string(data) != "foo2" {
		t.Fatal("Get() after Set()", err, string(data))
	}
}

func closeIt(v interface{}) {
	if c, ok := v.(io.Closer); ok {
		c.Close()
	}
}
