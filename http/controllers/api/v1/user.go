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

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/http/middlewares"
	"github.com/primasio/wormhole/http/oauth"
	"github.com/primasio/wormhole/http/token"
	"github.com/primasio/wormhole/models"
	"log"
)

type UserController struct{}

type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Remember string `form:"remember" json:"remember"`
}

type RegisterForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Nickname string `form:"nickname" json:"nickname" binding:"required"`
}

func (ctrl *UserController) Create(c *gin.Context) {

	var form RegisterForm

	if err := c.ShouldBind(&form); err != nil {
		Error(err.Error(), c)
	} else {
		dbi := db.GetDb()

		// Check username uniqueness

		exist := &models.User{}
		exist.Username = form.Username

		dbi.Where(&exist).First(&exist)

		if exist.ID != 0 {
			Error("Username exists", c)
			return
		}

		// Save user to db

		user := &models.User{}
		user.Username = form.Username
		user.Password = form.Password
		user.Nickname = form.Nickname

		if err := user.SetUniqueID(dbi); err != nil {
			ErrorServer(err, c)
		}

		dbi2 := dbi.Create(&user)

		if dbi2.Error != nil {
			ErrorServer(dbi2.Error, c)
		}

		Success(user, c)
	}
}

func (ctrl *UserController) Get(c *gin.Context) {

	userId, _ := c.Get(middlewares.AuthorizedUserId)

	user := &models.User{}
	user.ID = userId.(uint)

	dbi := db.GetDb()
	dbi.First(&user)

	if user.Username == "" {
		Error("User not found", c)
		return
	}

	Success(user, c)
}

func (ctrl *UserController) Auth(c *gin.Context) {

	var login LoginForm

	if err := c.ShouldBind(&login); err != nil {
		Error(err.Error(), c)
	} else {

		user := &models.User{Username: login.Username}

		dbi := db.GetDb()
		dbi.First(&user)

		if user.ID == 0 {
			ErrorUnauthorized("User not found", c)
			return
		}

		if !user.VerifyPassword(login.Password) {
			ErrorUnauthorized("Incorrect password", c)
		} else {

			// Login success, generate token
			accessToken, err := token.IssueToken(user.ID, login.Remember == "")

			if err != nil {
				ErrorServer(err, c)
			}

			Success(accessToken, c)
		}
	}
}

func (ctrl *UserController) GoogleAuth(c *gin.Context) {

	redirectUrl, err := oauth.HandleGoogleAuth()

	if err != nil {
		log.Println(err)
		Error(err.Error(), c)
	}

	c.Redirect(301, redirectUrl)
}

func (ctrl *UserController) GoogleAuthCallback(c *gin.Context) {

	token := c.Param("token")

	oauth.HandleGoogleAuthCallback(token)
}
