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

import (
	"github.com/gin-gonic/gin"
	"github.com/primasio/wormhole/http/controllers/api/v1"
	"github.com/primasio/wormhole/http/middlewares"
)

func NewRouter() *gin.Engine {

	gin.DisableConsoleColor()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(SetResponseHeader())

	v1g := router.Group("v1")
	{
		// OAuth 2.0 endpoints

		oauthCtrl := new(v1.OAuthController)

		oauthGroup := v1g.Group("oauth")
		{
			oauthGroup.GET("google", oauthCtrl.GoogleAuth)
			oauthGroup.POST("callback/google", oauthCtrl.GoogleAuthCallback)
		}

		// User endpoints

		userCtrl := new(v1.UserController)

		userGroup := v1g.Group("users")
		{
			userGroup.POST("/auth", userCtrl.Auth)
			userGroup.POST("", userCtrl.Create)

			userGroup.Use(middlewares.AuthMiddleware())
			{
				userGroup.GET("", userCtrl.Get)
			}
		}

		// Article endpoints

		articleCtrl := new(v1.ArticleController)

		articleGroup := v1g.Group("articles")
		{
			articleGroup.Use(middlewares.AuthMiddleware())
			{
				articleGroup.GET("/:article_id", articleCtrl.Get)
				articleGroup.POST("", articleCtrl.Publish)
			}
		}
	}

	return router
}
