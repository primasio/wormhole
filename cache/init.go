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
	"github.com/primasio/wormhole/config"
	"time"
)

var cacheStore CacheStore
var cacheType string

func InitCache() error {
	c := config.GetConfig()

	cacheType = c.GetString("cache.type")

	if cacheType == "memory" {
		cacheStore = NewInMemoryStore(time.Second)
	} else if cacheType == "redis" {

		host := c.GetString("cache.host")
		port := c.GetString("cache.port")
		password := c.GetString("cache.password")

		// Use our own redis cache since the original version is poorly written
		cacheStore = NewRedisCache(host+":"+port, password, time.Second)
	} else {
		return errors.New("unrecognized cache type")
	}

	return nil
}

func GetCache() CacheStore {
	return cacheStore
}

func GetCacheType() string {
	return cacheType
}
