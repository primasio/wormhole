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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/primasio/wormhole/cache"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/http/middlewares"
	"github.com/primasio/wormhole/models"
	"log"
	"strconv"
)

type UserController struct{}

type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Remember string `form:"remember" json:"remember"`
}

func (ctrl *UserController) Create(c *gin.Context) {

	var user models.User

	if err := c.ShouldBind(&user); err != nil {
		Error(err.Error(), c)
	} else {

		dbi := db.GetDb()

		// Check username uniqueness

		exist := &models.User{}
		exist.Username = user.Username

		dbi.First(&exist)

		if exist.ID != 0 {
			Error("Username exists", c)
			return
		}

		// Save user to db

		dbi2 := dbi.Create(&user)

		if dbi2.Error != nil {
			ErrorServer(dbi2.Error, c)
		}

		Success(user, c)
	}
}

func (ctrl *UserController) Get(c *gin.Context) {

	userId, exists := c.Get(middlewares.AuthorizedUserId)

	if !exists {
		Error("User not found", c)
		return
	}

	user := &models.User{}
	userIdInt, err2 := strconv.Atoi(fmt.Sprintf("%s", userId))

	if err2 != nil {
		log.Fatal(err2)
		Error("User not found", c)
		return
	}

	user.ID = uint(userIdInt)

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
			err, token := cache.NewSessionKey()

			if err != nil {
				log.Fatal(err)
				c.AbortWithStatus(500)
				return
			}

			userIdStr := fmt.Sprint(user.ID)
			cache.SessionSet(token, userIdStr, login.Remember == "")

			tokenStruct := make(map[string]string)
			tokenStruct["token"] = token

			Success(tokenStruct, c)
		}
	}
}
