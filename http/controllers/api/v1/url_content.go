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
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/http/middlewares"
	"github.com/primasio/wormhole/models"
)

type URLContentController struct{}

type URLContentForm struct {
	URL      string `form:"url" json:"url" binding:"required"`
	Content  string `form:"content" json:"content" binding:"required"`
	Abstract string `form:"abstract" json:"abstract" binding:"required"`
}

func (ctrl *URLContentController) Create(c *gin.Context) {
	var form URLContentForm

	if err := c.ShouldBind(&form); err != nil {
		Error(err.Error(), c)
	} else {
		dbi := db.GetDb()

		// Check URL uniqueness
		exist := &models.URLContent{}
		exist.URL = form.URL

		dbi.Where(&exist).First(&exist)

		if exist.ID != 0 {
			Error("URL exists", c)
			return
		}

		// Save url to db
		urlContent := &models.URLContent{}

		urlContent.URL = form.URL
		urlContent.Content = form.Content
		urlContent.Abstract = form.Abstract

		userId, _ := c.Get(middlewares.AuthorizedUserId)
		urlContent.UserId = userId.(uint)

		dbi.Create(&urlContent)

		Success(urlContent, c)
	}
}

func (ctrl *URLContentController) List(c *gin.Context) {

}

func (ctrl *URLContentController) Get(c *gin.Context) {

	url := c.Query("url")

	if url == "" {
		ErrorNotFound(errors.New("url not found"), c)
		return
	}

	urlContent := &models.URLContent{}

	urlContent.URL = url

	dbi := db.GetDb()
	dbi.Where(&urlContent).First(&urlContent)

	if urlContent.ID == 0 {
		ErrorNotFound(errors.New("url not found"), c)
		return
	}

	Success(urlContent, c)
}

func (ctrl *URLContentController) Vote(c *gin.Context) {

}
