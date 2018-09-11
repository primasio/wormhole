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

func Migration20180911() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "201809111335",
			Migrate: func(tx *gorm.DB) error {

				// it's a good practice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time

				type BaseModel struct {
					ID        uint `gorm:"primary_key" json:"-"`
					CreatedAt uint `json:"created_at"`
					UpdatedAt uint `json:"updated_at"`
				}

				type Article struct {
					BaseModel

					UserID   uint   `gorm:"index" json:"-"`
					Title    string `gorm:"type:text" form:"title" json:"title" binding:"required"`
					Abstract string `gorm:"type:text" json:"abstract"`
					Content  string `gorm:"type:longtext" form:"content" json:"content" binding:"required"`
					Language string `gorm:"column:lang;size:64" json:"language"`

					ContentId  string `gorm:"type:varchar(128);unique_index" json:"content_id"`
					ContentDNA string `gorm:"type:varchar(128);unique_index" json:"content_dna"`
				}

				type Domain struct {
					BaseModel
					UserID  uint   `json:"-"`
					Domain  string `gorm:"type:text" json:"domain"`
					Title   string `gorm:"type:text" json:"title"`
					HashKey string `gorm:"type:varchar(128);unique_index" json:"-"`

					IsActive bool `gorm:"default:false" json:"is_active"`
					Votes    uint `gorm:"default:1" json:"votes"`
				}

				type DomainVote struct {
					BaseModel
					UserID   uint
					DomainID uint
				}

				type URLContent struct {
					BaseModel
					UserID  uint   `json:"-"`
					URL     string `gorm:"type:text" json:"url"`
					HashKey string `gorm:"type:varchar(128);unique_index" json:"-"`

					TotalComment uint `gorm:"default:1" json:"total_comment"`
				}

				type URLContentComment struct {
					BaseModel

					UniqueID     string `gorm:"type:varchar(128);unique_index" json:"id"`
					UserID       uint   `json:"-"`
					URLContentId uint   `json:"-"`
					Content      string `gorm:"type:longtext" json:"content"`

					IsDeleted bool `gorm:"default:false" json:"is_deleted"`
				}

				if err := tx.AutoMigrate(&Article{}).Error; err != nil {
					return err
				}

				if err := tx.AutoMigrate(&Domain{}).Error; err != nil {
					return err
				}

				if err := tx.AutoMigrate(&DomainVote{}).Error; err != nil {
					return err
				}

				if err := tx.AutoMigrate(&URLContent{}).Error; err != nil {
					return err
				}

				if err := tx.AutoMigrate(&URLContentComment{}).Error; err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.DropTable("articles").Error; err != nil {
					return err
				}
				if err := tx.DropTable("domains").Error; err != nil {
					return err
				}
				if err := tx.DropTable("domain_votes").Error; err != nil {
					return err
				}
				if err := tx.DropTable("url_contents").Error; err != nil {
					return err
				}
				if err := tx.DropTable("url_content_comments").Error; err != nil {
					return err
				}
				return nil
			},
		},
	}
}
