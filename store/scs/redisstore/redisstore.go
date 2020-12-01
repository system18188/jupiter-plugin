package redisstore

import (
	"errors"
	"sync"
	"time"

	"github.com/douyu/jupiter/pkg/client/redis"
)

var (
	errorCommit = errors.New("scs Commit Set Redis Error")
)

// RedisStore represents the session store.
type RedisStore struct {
	redis *redis.Redis
	*sync.RWMutex
	prefix string
}

// New returns a new RedisStore instance. The pool parameter should be a pointer
// to a redigo connection pool. See https://godoc.org/github.com/gomodule/redigo/redis#Pool.
func New(redis *redis.Redis) *RedisStore {
	return NewWithPrefix(redis, "scs:session:")
}

// NewWithPrefix returns a new RedisStore instance. The pool parameter should be a pointer
// to a redigo connection pool. The prefix parameter controls the Redis key
// prefix, which can be used to avoid naming clashes if necessary.
func NewWithPrefix(redis *redis.Redis, prefix string) *RedisStore {
	return &RedisStore{
		redis:   redis,
		prefix:  prefix,
		RWMutex: new(sync.RWMutex),
	}
}

// Find returns the data for a given session token from the RedisStore instance.
// If the session token is not found or is expired, the returned exists flag
// will be set to false.
func (r *RedisStore) Find(token string) (b []byte, exists bool, err error) {
	r.RWMutex.RLock()
	defer r.RWMutex.RUnlock()
	b, err = r.redis.GetRaw(r.prefix + token)
	if err != nil {
		return nil, false, err
	}
	return b, true, nil
}

// Commit adds a session token and data to the RedisStore instance with the
// given expiry time. If the session token already exists then the data and
// expiry time are updated.
func (r *RedisStore) Commit(token string, b []byte, expiry time.Time) error {
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()
	if !r.redis.Set(r.prefix+token, b, makeMillisecondTimestamp(expiry)) {
		return errorCommit
	}
	return nil
}

// Delete removes a session token and corresponding data from the RedisStore
// instance.
func (r *RedisStore) Delete(token string) error {
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()
	r.redis.Del(r.prefix+token)
	return nil
}

func makeMillisecondTimestamp(t time.Time) time.Duration {
	return time.Duration(t.UnixNano()) / (time.Millisecond / time.Nanosecond)
}
