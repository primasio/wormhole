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
	"strconv"
)

type URLContentController struct{}

type URLContentForm struct {
	URL      string `form:"url" json:"url" binding:"required"`
	Title    string `form:"title" json:"title" binding:"required"`
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
		urlContent.Title = form.Title
		urlContent.Content = form.Content
		urlContent.Abstract = form.Abstract

		userId, _ := c.Get(middlewares.AuthorizedUserId)
		urlContent.UserId = userId.(uint)

		dbi.Create(&urlContent)

		Success(urlContent, c)
	}
}

func (ctrl *URLContentController) List(c *gin.Context) {

	urlType := c.Query("type")

	pageSize := 20
	page, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		page = 0
	}

	offsetNum := page * pageSize

	var urlList []models.URLContent

	dbi := db.GetDb()
	query := dbi.Select("id, url, title, abstract, votes, is_active, total_comment, created_at, updated_at")
	query = query.Where("is_active = ?", urlType != "voting")
	query = query.Order("created_at DESC").Offset(offsetNum).Limit(pageSize)

	if err := query.Find(&urlList).Error; err != nil {
		ErrorServer(err, c)
		return
	}

	Success(urlList, c)
}

func (ctrl *URLContentController) Get(c *gin.Context) {

	url := c.Query("url")

	err, urlContent := findURL(url)

	if err != nil {
		ErrorNotFound(err, c)
		return
	}

	Success(urlContent, c)
}

func (ctrl *URLContentController) Vote(c *gin.Context) {

	url := c.Query("url")

	err, urlContent := findURL(url)

	if err != nil {
		ErrorNotFound(err, c)
		return
	}

	if urlContent.IsActive {
		Error("url is already active", c)
		return
	}

	userId, _ := c.Get(middlewares.AuthorizedUserId)
	userIdNum := userId.(uint)

	if urlContent.UserId == userId {
		Error("user already voted", c)
		return
	}

	// Update vote should be executed on a locked object

	lockedURLContent := &models.URLContent{}

	tx := db.GetDb().Begin()

	// If SQLite is used, FOR UPDATE is not supported
	// Then there is an error of concurrent votes count

	sql := "SELECT * FROM url_contents WHERE id = ?"

	if db.GetDbType() != "sqlite3" {
		sql = sql + " FOR UPDATE"
	}

	if err := tx.Raw(sql, urlContent.ID).Scan(&lockedURLContent).Error; err != nil {
		tx.Rollback()
		ErrorServer(err, c)
		return
	}

	if lockedURLContent.ID == 0 {
		tx.Rollback()
		ErrorServer(errors.New("error lock url_content"), c)
		return
	}

	// Check user vote status
	// this should be performed after the locking of url_content
	// to avoid race condition of concurrent voting from the same user

	vote := &models.URLContentVote{
		UserId:       userIdNum,
		URLContentID: lockedURLContent.ID,
	}

	tx.Where(&vote).First(&vote)

	if vote.ID != 0 {
		tx.Rollback()
		Error("user already voted", c)
		return
	}

	lockedURLContent.Votes++

	if err := tx.Save(&lockedURLContent).Error; err != nil {
		tx.Rollback()
		ErrorServer(err, c)
		return
	}

	if err := tx.Create(vote).Error; err != nil {
		tx.Rollback()
		ErrorServer(err, c)
		return
	}

	tx.Commit()
	Success(lockedURLContent, c)
}

func findURL(url string) (error, *models.URLContent) {
	if url == "" {
		return errors.New("url is empty"), nil
	}

	urlContent := &models.URLContent{}

	urlContent.URL = url

	dbi := db.GetDb()
	dbi.Where(&urlContent).First(&urlContent)

	if urlContent.ID == 0 {
		return errors.New("url not found"), nil
	}

	return nil, urlContent
}
