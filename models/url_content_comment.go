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

package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/primasio/wormhole/util"
)

type URLContentComment struct {
	BaseModel

	UniqueID     string `gorm:"unique_index" json:"id"`
	UserId       uint   `json:"-"`
	URLContentId uint   `json:"-"`
	Content      string `gorm:"type:longtext" json:"content"`

	User User `gorm:"save_associations:false" sql:"-" json:"user"`

	IsDeleted bool `gorm:"default:false" json:"is_deleted"`
}

func (comment *URLContentComment) SetUniqueID(db *gorm.DB) error {
	var counter = 0

	for {
		counter = counter + 1
		uid := util.RandStringUppercase(8)

		check := &URLContentComment{UniqueID: uid}

		db.Where(&check).First(&check)

		if check.ID == 0 {
			comment.UniqueID = uid
			return nil
		}

		if counter >= 5 {
			// This is unlikely to happen
			// Must be error from other parts
			return errors.New("too many iterations while generating new session key")
		}
	}
}
