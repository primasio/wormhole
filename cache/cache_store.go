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
)

const (
	DEFAULT = time.Duration(0)
	FOREVER = time.Duration(-1)
)

var (
	PageCachePrefix = "gincontrib.page.cache"
	ErrCacheMiss    = errors.New("cache: key not found")
	ErrNotStored    = errors.New("cache: not stored")
	ErrNotSupport   = errors.New("cache: not support")
)

// CacheStore is the interface of a cache backend
type CacheStore interface {
	// Get retrieves an item from the cache. Returns the item or nil, and a bool indicating
	// whether the key was found.
	Get(key string, value interface{}) error

	// Set sets an item to the cache, replacing any existing item.
	Set(key string, value interface{}, expire time.Duration) error

	// Add adds an item to the cache only if an item doesn't already exist for the given
	// key, or if the existing item has expired. Returns an error otherwise.
	Add(key string, value interface{}, expire time.Duration) error

	// Replace sets a new value for the cache key only if it already exists. Returns an
	// error if it does not.
	Replace(key string, data interface{}, expire time.Duration) error

	// Delete removes an item from the cache. Does nothing if the key is not in the cache.
	Delete(key string) error

	// Increment increments a real number, and returns error if the value is not real
	Increment(key string, data int64) (int64, error)

	// Decrement decrements a real number, and returns error if the value is not real
	Decrement(key string, data int64) (int64, error)

	// Flush seletes all items from the cache.
	Flush() error

	// Expire a key
	Expire(key string, expire time.Duration) (bool, error)
}
