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

package service

import (
	"sync"

	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/models"
)

var ucc *URLContentComment
var uccOnce sync.Once

type URLContentComment struct{}

func GetURLContentComment() *URLContentComment {
	uccOnce.Do(func() {
		ucc = &URLContentComment{}
	})

	return ucc
}

func (s *URLContentComment) ListWithVote(userID uint, urlContent *models.URLContent, page, pageSize, offsetNum int) interface{} {
	type ScanItem struct {
		models.BaseModel
		UniqueID         string
		Content          string
		CommentUpVotes   uint
		CommentDownVotes uint
		IsDeleted        bool

		UserUniqueID         string
		UserAvatarURL        string
		UserNickname         string
		UserIntegration      int64
		UserCommentUpVotes   uint
		UserCommentDownVotes uint
		UserBalance          string
		UserCreatedAt        uint
		UserUpdatedAt        uint

		Like string
	}

	type User struct {
		CreatedAt        uint   `json:"created_at"`
		UpdatedAt        uint   `json:"updated_at"`
		ID               string `json:"id"`
		Nickname         string `json:"nickname"`
		AvatarURL        string `json:"avatar_url"`
		Integration      int64  `json:"integration"`
		CommentUpVotes   uint   `json:"comment_up_votes"`
		CommentDownVotes uint   `json:"comment_down_votes"`
		Balance          string `json:"balance"`
	}

	type ResultItem struct {
		CreatedAt        uint   `json:"created_at"`
		UpdatedAt        uint   `json:"updated_at"`
		ID               string `json:"id"`
		Content          string `json:"content"`
		CommentUpVotes   uint   `json:"comment_up_votes"`
		CommentDownVotes uint   `json:"comment_down_votes"`
		User             User   `json:"user"`
		Like             string `json:"like"`
		IsDeleted        bool   `json:"is_deleted"`
	}

	items := make([]ResultItem, 0)

	dbi := db.GetDb()
	query := dbi.Table("url_content_comments")
	rows, _ := query.Order("url_content_comments.created_at DESC").
		Select("users.integration as user_integration, users.comment_up_votes as user_comment_up_votes, users.comment_down_votes as user_comment_down_votes, users.balance as user_balance,users.created_at as user_created_at, users.updated_at as users_updated_at,users.avatar_url as user_avatar_url, users.nickname as user_nickname, users.unique_id as user_unique_id, url_content_comments.*, url_content_comment_votes.like").
		Joins("left join users on url_content_comments.user_id = users.id").
		Joins("left join url_content_comment_votes on url_content_comment_votes.url_content_comment_id = url_content_comments.id and url_content_comment_votes.user_id = ?", userID).
		Where("url_content_comments.url_content_id = ? AND url_content_comments.is_deleted = 0", urlContent.ID).
		Offset(offsetNum).
		Limit(pageSize).Rows()

	defer rows.Close()

	for rows.Next() {
		v := &ScanItem{}
		dbi.ScanRows(rows, &v)

		result := ResultItem{
			ID:               v.UniqueID,
			Content:          v.Content,
			CommentUpVotes:   v.CommentUpVotes,
			CommentDownVotes: v.CommentDownVotes,
			IsDeleted:        v.IsDeleted,
			CreatedAt:        v.CreatedAt,
			UpdatedAt:        v.UpdatedAt,
			Like:             v.Like,

			User: User{
				CreatedAt: v.UserCreatedAt, UpdatedAt: v.UserUpdatedAt, ID: v.UserUniqueID, Nickname: v.UserNickname, AvatarURL: v.UserAvatarURL,
				Integration: v.UserIntegration, CommentUpVotes: v.UserCommentUpVotes, CommentDownVotes: v.UserCommentDownVotes, Balance: v.UserBalance,
			},
		}

		items = append(items, result)
	}

	return items
}
