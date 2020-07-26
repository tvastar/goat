package sessions_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/tvastar/goat"
	"github.com/tvastar/goat/sessions"
)

func TestRedis(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("miniredis start", err)
	}
	defer s.Close()

	var r goat.SessionStore = &sessions.Redis{
		Options: redis.Options{Addr: s.Addr()},
		TTL:     time.Second,
	}

	ctx := context.Background()
	data, err := r.Get(ctx, "foo", "nonce")
	if err == nil || data != nil {
		t.Fatal("Unexpected initial fetch", data, err)
	}

	err = r.Set(ctx, "foo", "nonce", []byte{0, 1, 2, 3})
	if err != nil {
		t.Fatal("Unexpected set failure", err)
	}

	data, err = r.Get(ctx, "foo", "nonce")
	if err != nil || string(data) != string([]byte{0, 1, 2, 3}) {
		t.Fatal("Unexpected get failure", data, err)
	}

	data, err = r.Get(ctx, "foo", "nonce")
	if err == nil || data != nil {
		t.Fatal("Unexpected second fetch", data, err)
	}
}
