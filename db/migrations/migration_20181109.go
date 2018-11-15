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

func Migration20181109() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "201811091130",
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
					Balance   string `json:"balance"`
				}

				if err := tx.AutoMigrate(&User{}).Error; err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	}
}
