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
		tx := db.GetDb().Begin()

		// Check URL content
		err, lockedUrlContent := models.GetURLContentByURL(form.URL, tx, true)

		if err != nil {
			tx.Rollback()
			ErrorServer(err, c)
			return
		}

		if lockedUrlContent == nil {
			tx.Rollback()
			ErrorNotFound(errors.New("url not found"), c)
			return
		}

		userId, _ := c.Get(middlewares.AuthorizedUserId)

		// Create comment

		comment := models.URLContentComment{}
		comment.UserId = userId.(uint)
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

		// Update comment count
		lockedUrlContent.TotalComment++

		if err = tx.Save(&lockedUrlContent).Error; err != nil {
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

	comment := &models.URLContentComment{}
	comment.UniqueID = commentId

	dbi := db.GetDb()
	dbi.Where(comment).First(&comment)

	if comment.ID == 0 {
		ErrorNotFound(errors.New("comment not found"), c)
		return
	}

	comment.IsDeleted = true
	dbi.Save(comment)

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

		if urlContent == nil {
			ErrorNotFound(errors.New("url not found"), c)
			return
		}

		var commentList []models.URLContentComment

		query := dbi.Where("url_content_id = ? AND is_deleted = 0", urlContent.ID)
		query.Order("created_at DESC").Offset(offsetNum).Limit(pageSize).Find(&commentList)

		Success(commentList, c)
	}
}
