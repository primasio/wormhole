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
	"github.com/gin-contrib/cache/persistence"
	"github.com/primasio/wormhole/config"
	"time"
)

var cacheStore persistence.CacheStore

func InitCache() {
	c := config.GetConfig()

	cacheType := c.GetString("cache.type")

	if cacheType == "memory" {
		cacheStore = persistence.NewInMemoryStore(time.Second)
	} else {

		host := c.GetString("cache.host")
		port := c.GetString("cache.port")
		password := c.GetString("cache.password")

		cacheStore = persistence.NewRedisCache(host+":"+port, password, time.Second)
	}
}

func GetCache() persistence.CacheStore {
	return cacheStore
}
