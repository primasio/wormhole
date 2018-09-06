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

package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/primasio/wormhole/cache"
	"net/http"
	"strconv"
	"time"
)

const AuthorizedUserId = "UserId"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		reqToken := c.Request.Header.Get("Authorization")

		if reqToken == "" {
			c.AbortWithStatus(401)
			return
		}

		// Check token validity

		if err, userId := cache.SessionGet(reqToken); err != nil {

			glog.Error("token not exist", err)
			c.AbortWithStatus(500)

		} else {

			if userId == "" {
				c.AbortWithStatus(401)
			} else {

				userIdNum, err := strconv.Atoi(userId)

				if err != nil {
					glog.Error(err)
					c.AbortWithStatus(500)
					return
				}

				c.Set(AuthorizedUserId, uint(userIdNum))

				// User account based access rate limit

				err, reached := rateLimitReached(userId)

				if err != nil {
					glog.Error(err)
					c.AbortWithStatus(http.StatusInternalServerError)
				} else {
					if reached {
						c.AbortWithStatus(http.StatusBadRequest)
					} else {
						c.Next()
					}
				}
			}
		}
	}
}

func rateLimitReached(userId string) (error, bool) {

	cacheType := cache.GetCacheType()

	if cacheType == "memory" {
		return nil, false
	}

	// API access for a single user is limited to 10 times per minute

	currentMinute := int(time.Now().Unix() / 60)

	slotId := "auth_rate_limit_" + userId + "_" + strconv.Itoa(currentMinute)

	cache := cache.GetCache()

	count, err := cache.Increment(slotId, 1)

	if err != nil {
		return err, false
	}

	if count >= 10 {
		return nil, true
	} else {
		cache.Expire(slotId, time.Minute)
	}

	return nil, false
}
