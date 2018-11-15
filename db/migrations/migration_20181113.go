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

package migrations

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func Migration20181113() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "201811131031",
			Migrate: func(tx *gorm.DB) error {

				type BaseModel struct {
					ID        uint `gorm:"primary_key" json:"-"`
					CreatedAt uint `json:"created_at"`
					UpdatedAt uint `json:"updated_at"`
				}

				type User struct {
					BaseModel
					UniqueID  string `json:"id" gorm:"type:varchar(128);unique_index"`
					Username  string `json:"-" gorm:"type:varchar(128);index"`
					Password  string `json:"-"`
					Salt      string `json:"-"`
					Nickname  string `json:"nickname"`
					AvatarURL string `json:"avatar_url"`

					Integration      int64 `json:"integration" gorm:"type:INT(18);default:0"`
					CommentUpVotes   uint  `json:"comment_up_votes" gorm:"type:INT(11);default:0"`
					CommentDownVotes uint  `json:"comment_down_votes" gorm:"type:INT(11);default:0"`

					Balance string `json:"balance"`
				}

				type URLContentComment struct {
					BaseModel

					UniqueID         string `gorm:"type:varchar(128);unique_index" json:"id"`
					UserID           uint   `json:"-"`
					URLContentId     uint   `json:"-"`
					Content          string `gorm:"type:longtext" json:"content"`
					CommentUpVotes   uint   `json:"comment_up_votes" gorm:"default:0"`
					CommentDownVotes uint   `json:"comment_down_votes" gorm:"default:0"`
					IsDeleted        bool   `gorm:"default:false" json:"is_deleted"`
				}

				type URLContentCommentVote struct {
					BaseModel
					UniqueID            string `json:"id" gorm:"type:varchar(128);unique_index"`
					UserID              uint   `json:"-"`
					URLContentCommentID uint   `json:"-"`
					Like                bool   `json:"like"`
				}

				type IntegrationHistory struct {
					BaseModel
					UniqueID    string `json:"id" gorm:"type:varchar(128);unique_index"`
					UserID      uint   `json:"-"`
					Integration int64  `json:"integration"`
					Description string `json:"description"`
					Data        string `json:"-"`
				}

				type RegisterIntegrationWorkerInfo struct {
					BaseModel
					LastDoneUserID uint `json:"-"`
				}

				if err := tx.AutoMigrate(&User{}).Error; err != nil {
					return err
				}

				if err := tx.AutoMigrate(&URLContentComment{}).Error; err != nil {
					return err
				}

				if err := tx.AutoMigrate(&URLContentCommentVote{}).Error; err != nil {
					return err
				}

				if err := tx.AutoMigrate(&IntegrationHistory{}).Error; err != nil {
					return err
				}

				if err := tx.AutoMigrate(&RegisterIntegrationWorkerInfo{}).Error; err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.DropTable("url_content_comment_votes").Error; err != nil {
					return err
				}
				if err := tx.DropTable("url_content_comments").Error; err != nil {
					return err
				}
				if err := tx.DropTable("integration_histories").Error; err != nil {
					return err
				}
				if err := tx.DropTable("register_integration_worker_infos").Error; err != nil {
					return err
				}
				return nil
			},
		},
	}
}
