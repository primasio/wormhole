/*
 * Copyright 2018 Primas Lab Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cache

import (
	"errors"
	"time"

	"github.com/dpordomingo/go-gingonic-cache/utils"
	"github.com/garyburd/redigo/redis"
)

type RedisStore struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

// NewRedisCache returns a RedisStore
// until redigo supports sharding/clustering, only one host will be in hostList
func NewRedisCache(host string, password string, defaultExpiration time.Duration) *RedisStore {
	var pool = &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			// the redis protocol should probably be made sett-able
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			} else {
				// check with PING
				if _, err := c.Do("PING"); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		// custom connection test method
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				return err
			}
			return nil
		},
	}
	return &RedisStore{pool, defaultExpiration}
}

// Set (see CacheStore interface)
func (c *RedisStore) Set(key string, value interface{}, expires time.Duration) error {
	return c.invoke(c.pool.Get().Do, key, value, expires)
}

// Add (see CacheStore interface)
func (c *RedisStore) Add(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	if exists(conn, key) {
		return ErrNotStored
	}
	return c.invoke(conn.Do, key, value, expires)
}

// Replace (see CacheStore interface)
func (c *RedisStore) Replace(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	if !exists(conn, key) {
		return ErrNotStored
	}
	err := c.invoke(conn.Do, key, value, expires)
	if value == nil {
		return ErrNotStored
	}

	return err
}

// Get (see CacheStore interface)
func (c *RedisStore) Get(key string, ptrValue interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	raw, err := conn.Do("GET", key)
	if raw == nil {
		return ErrCacheMiss
	}
	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}
	return utils.Deserialize(item, ptrValue)
}

func exists(conn redis.Conn, key string) bool {
	retval, _ := redis.Bool(conn.Do("EXISTS", key))
	return retval
}

// Delete (see CacheStore interface)
func (c *RedisStore) Delete(key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	if !exists(conn, key) {
		return ErrCacheMiss
	}
	_, err := conn.Do("DEL", key)
	return err
}

// Increment (see CacheStore interface)
func (c *RedisStore) Increment(key string, delta int64) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	ret, err := conn.Do("INCRBY", key, delta)
	return redis.Int64(ret, err)
}

// Decrement (see CacheStore interface)
func (c *RedisStore) Decrement(key string, delta int64) (newValue int64, err error) {
	conn := c.pool.Get()
	defer conn.Close()
	ret, err := conn.Do("DECRBY", key, delta)
	return redis.Int64(ret, err)
}

// Flush (see CacheStore interface)
func (c *RedisStore) Flush() error {
	return errors.New("not implemented")
}

func (c *RedisStore) Expire(key string, expires time.Duration) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	ret, err := conn.Do("EXPIRE", key, int32(expires/time.Second))

	return redis.Bool(ret, err)
}

func (c *RedisStore) invoke(f func(string, ...interface{}) (interface{}, error),
	key string, value interface{}, expires time.Duration) error {

	switch expires {
	case DEFAULT:
		expires = c.defaultExpiration
	case FOREVER:
		expires = time.Duration(0)
	}

	b, err := utils.Serialize(value)
	if err != nil {
		return err
	}
	conn := c.pool.Get()
	defer conn.Close()
	if expires > 0 {
		_, err := f("SETEX", key, int32(expires/time.Second), b)
		return err
	}

	_, err = f("SET", key, b)
	return err
}
