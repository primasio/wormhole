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
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/primasio/wormhole/config"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/http/captcha"
	"github.com/primasio/wormhole/http/middlewares"
	"github.com/primasio/wormhole/models"
	"github.com/primasio/wormhole/service"
)

type URLContentCommentController struct{}

type URLContentCommentForm struct {
	URL     string `form:"url" json:"url" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}

func (ctrl *URLContentCommentController) Create(c *gin.Context) {
	var form URLContentCommentForm

	if err := c.ShouldBind(&form); err != nil {
		Error(err.Error(), c)
	} else {

		err, _ := models.ExtractDomainFromURL(form.URL)

		if err != nil {
			ErrorServer(err, c)
			return
		}

		tx := db.GetDb().Begin()

		// Check URL content
		err, lockedUrlContent := models.GetURLContentByURL(form.URL, tx, true)

		if err != nil {
			tx.Rollback()
			ErrorServer(err, c)
			return
		}

		userId, _ := c.Get(middlewares.AuthorizedUserId)

		if lockedUrlContent == nil {
			// First time comment
			// Create the url content

			lockedUrlContent = &models.URLContent{}
			lockedUrlContent.UserID = userId.(uint)
			lockedUrlContent.URL = models.CleanURL(form.URL)
			lockedUrlContent.HashKey = models.GetURLHashKey(lockedUrlContent.URL)

			err = tx.Create(lockedUrlContent).Error

		} else {
			// Update comment count
			lockedUrlContent.TotalComment++

			err = tx.Save(&lockedUrlContent).Error
		}

		if err != nil {
			tx.Rollback()
			ErrorServer(err, c)
			return
		}

		// Create comment

		comment := models.URLContentComment{}
		comment.UserID = userId.(uint)
		comment.URLContentId = lockedUrlContent.ID
		comment.Content = form.Content

		if err := comment.SetUniqueID(tx); err != nil {
			tx.Rollback()
			ErrorServer(err, c)
			return
		}

		if err = tx.Create(&comment).Error; err != nil {
			tx.Rollback()
			ErrorServer(err, c)
			return
		}

		tx.Commit()
		Success(comment, c)
	}
}

func (ctrl *URLContentCommentController) Delete(c *gin.Context) {

	commentId := c.Param("comment_id")

	if commentId == "" {
		Error("missing comment id", c)
		return
	}

	if config.GetAppEnvironment() == config.AppEnvProduction {

		// Check captcha

		token := c.Query("token")

		if token == "" {
			Error("missing query param token", c)
			return
		}

		err, passed := captcha.VerifyRecaptchaToken(token)

		if err != nil {
			Error(err.Error(), c)
			return
		}

		if !passed {
			Error("captcha verification failed", c)
			return
		}
	}

	comment := &models.URLContentComment{}
	comment.UniqueID = commentId

	tx := db.GetDb().Begin()
	tx.Where(comment).First(&comment)

	if comment.ID == 0 {
		tx.Rollback()
		ErrorNotFound(errors.New("comment not found"), c)
		return
	}

	sql := "SELECT id, hash_key, total_comment FROM url_contents WHERE id = ?"

	if db.GetDbType() != db.SQLITE {
		sql = sql + " FOR UPDATE"
	}

	var urlContent models.URLContent

	tx.Raw(sql, comment.URLContentId).Scan(&urlContent)

	if urlContent.ID == 0 {
		tx.Rollback()
		ErrorNotFound(errors.New("url not found"), c)
		return
	}

	lockedComment := &models.URLContentComment{}
	lockedComment.ID = comment.ID

	tx.Where(lockedComment).First(&lockedComment)

	if lockedComment.CreatedAt == 0 {
		tx.Rollback()
		ErrorNotFound(errors.New("comment not found"), c)
		return
	}

	lockedComment.IsDeleted = true
	if err := tx.Save(lockedComment).Error; err != nil {
		tx.Rollback()
		ErrorServer(err, c)
		return
	}

	urlContent.TotalComment--

	if err := tx.Save(&urlContent).Error; err != nil {
		tx.Rollback()
		ErrorServer(err, c)
		return
	}

	tx.Commit()

	Success(nil, c)
}

func (ctrl *URLContentCommentController) List(c *gin.Context) {

	page, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		page = 0
	}

	pageSize := 20
	offsetNum := page * pageSize

	url := c.Query("url")

	if url == "" {
		Error("missing query param url", c)
		return
	}

	dbi := db.GetDb()
	if err, urlContent := models.GetURLContentByURL(url, dbi, false); err != nil {
		ErrorServer(err, c)
		return
	} else {

		commentList := make([]models.URLContentComment, 0)

		if urlContent != nil {
			query := dbi.Where("url_content_id = ? AND is_deleted = 0", urlContent.ID)
			query.Order("created_at DESC").Offset(offsetNum).Limit(pageSize).Preload("User").Find(&commentList)
		}

		Success(commentList, c)
	}
}

func (ctrl *URLContentCommentController) ListWithVote(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		page = 0
	}

	pageSize := 20
	offsetNum := page * pageSize

	url := c.Query("url")

	if url == "" {
		Error("missing query param url", c)
		return
	}

	dbi := db.GetDb()
	if err, urlContent := models.GetURLContentByURL(url, dbi, false); err != nil {
		ErrorServer(err, c)
		return
	} else {

		if urlContent == nil {
			Success(make([]interface{}, 0), c)
			return
		}

		userID, _ := c.Get(middlewares.AuthorizedUserId)
		items := service.GetURLContentComment().ListWithVote(userID.(uint), urlContent, page, pageSize, offsetNum)

		Success(items, c)
	}
}
