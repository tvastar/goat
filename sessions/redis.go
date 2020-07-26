// Package sessions implements a session store for goat.
package sessions

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis implements a session store on top of Redis
type Redis struct {
	Options redis.Options
	TTL     time.Duration
	c       *redis.Client
	sync.Once
}

// Get returns a key
func (r *Redis) Get(ctx context.Context, provider, nonce string) ([]byte, error) {
	client := r.init()

	var result *redis.StringCmd
	_, err := client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		result = client.Get(ctx, provider+"/"+nonce)
		client.Del(ctx, provider+"/"+nonce)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result.Bytes()

}

func (r *Redis) Set(ctx context.Context, provider, nonce string, state []byte) error {
	client := r.init()

	return client.Set(ctx, provider+"/"+nonce, state, r.TTL).Err()
}

func (r *Redis) init() *redis.Client {
	r.Once.Do(func() {
		r.c = redis.NewClient(&r.Options)
	})
	return r.c
}
