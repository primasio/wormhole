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
	"github.com/primasio/wormhole/db"

	"github.com/primasio/wormhole/http/middlewares"
	"github.com/primasio/wormhole/models"

	"github.com/gin-gonic/gin"
	"github.com/primasio/wormhole/service"
)

type URLContentCommentVoteController struct{}

type URLContentCommentVoteForm struct {
	Like bool `form:"like" json:"like" binding:"exists"`
}

func (ctrl *URLContentCommentVoteController) Create(c *gin.Context) {
	var form URLContentCommentVoteForm

	if err := c.ShouldBind(&form); err != nil {

		Error(err.Error(), c)

	} else {

		dbi := db.GetDb()
		user := &models.User{}
		comment := &models.URLContentComment{}
		userID, _ := c.Get(middlewares.AuthorizedUserId)

		commentID := c.Param("comment_id")
		if err := dbi.Where("unique_id = ?", commentID).First(comment).Error; err != nil {
			Error(err.Error(), c)
			return
		}

		if err := dbi.Where("id = ?", userID.(uint)).First(user).Error; err != nil {
			Error(err.Error(), c)
			return
		}

		if err := service.GetURLContentCommentVote().CreateVote(dbi, comment, user, form.Like); err != nil {
			Error(err.Error(), c)
			return
		}

		Success(nil, c)

	}
}

func (ctrl *URLContentCommentVoteController) Update(c *gin.Context) {
	var form URLContentCommentVoteForm

	if err := c.ShouldBind(&form); err != nil {

		Error(err.Error(), c)

	} else {

		dbi := db.GetDb()
		user := &models.User{}
		commentID := c.Param("comment_id")
		comment := &models.URLContentComment{}
		userID, _ := c.Get(middlewares.AuthorizedUserId)

		if err := dbi.Where("unique_id = ?", commentID).First(comment).Error; err != nil {
			Error(err.Error(), c)
			return
		}

		if err := dbi.Where("id = ?", userID.(uint)).First(user).Error; err != nil {
			Error(err.Error(), c)
			return
		}

		if err := service.GetURLContentCommentVote().UpdateVote(dbi, comment, user, form.Like); err != nil {
			Error(err.Error(), c)
			return
		}

		Success(nil, c)

	}
}

func (ctrl *URLContentCommentVoteController) Delete(c *gin.Context) {

	dbi := db.GetDb()
	user := &models.User{}
	commentID := c.Param("comment_id")
	comment := &models.URLContentComment{}
	userID, _ := c.Get(middlewares.AuthorizedUserId)

	if err := dbi.Where("unique_id = ?", commentID).First(comment).Error; err != nil {
		Error(err.Error(), c)
		return
	}

	if err := dbi.Where("id = ?", userID.(uint)).First(user).Error; err != nil {
		Error(err.Error(), c)
		return
	}

	if err := service.GetURLContentCommentVote().CancelVote(dbi, comment, user); err != nil {
		Error(err.Error(), c)
		return
	}

	vote := &models.URLContentCommentVote{UserID: user.ID, URLContentCommentID: comment.ID}
	vote.SetUniqueID()

	Success(nil, c)

}
