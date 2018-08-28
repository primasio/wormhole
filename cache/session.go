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
	"github.com/primasio/wormhole/util"
	"time"
)

func NewSessionKey() (error, string) {

	var counter = 0

	for {
		counter = counter + 1
		key := util.RandString(32)
		err, check := SessionGet(key)

		if err != nil {
			return err, ""
		}

		if check == "" {
			return nil, key
		}

		if counter >= 5 {
			// This is unlikely to happen
			// Must be error from other parts
			return errors.New("too many iterations while generating new session key"), ""
		}
	}
}

func SessionSet(token, userId string, expires bool) error {
	store := GetCache()

	duration := time.Hour * 24 * 30

	if expires {
		duration = time.Hour * 2
	}

	return store.Set("wormhole_session_"+token, userId, duration)
}

func SessionGet(token string) (err error, userId string) {

	store := GetCache()

	var userIdStore string

	if err := store.Get("wormhole_session_"+token, &userIdStore); err != nil {
		return err, ""
	}

	return nil, userIdStore
}
