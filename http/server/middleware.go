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

package server

import "github.com/gin-gonic/gin"

var allowedOrigins map[string]bool

//SetResponseHeader set response header for all requests
func SetResponseHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqOrigin := c.Request.Header.Get("Origin")
		//log.Println("****************reqOrigin:", reqOrigin)

		if _, haveOrigin := allowedOrigins[reqOrigin]; haveOrigin {
			//c.Request.Close = true
			c.Header("Access-Control-Allow-Origin", reqOrigin)
			c.Header("Connection", "keep-alive")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
	}
}
