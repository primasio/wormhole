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
	"crypto/sha1"
	"encoding/hex"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/primasio/wormhole/db"
)

type URLContent struct {
	BaseModel
	UserID  uint   `json:"-"`
	URL     string `gorm:"type:text" json:"url"`
	HashKey string `gorm:"type:varchar(128);unique_index" json:"-"`

	TotalComment uint `gorm:"default:1" json:"total_comment"`
}

func CleanURL(url string) string {
	//TODO: Remove trailing slash, remove hash, etc.
	return url
}

func GetURLHashKey(url string) string {

	sumBytes := sha1.Sum([]byte(url))

	return hex.EncodeToString(sumBytes[:])
}

func GetURLContentCount(dbi *gorm.DB) (uint, error) {
	count := 0
	err := dbi.Model(&URLContent{}).Select("id").Count(&count).Error
	return uint(count), err
}

func GetURLContentByURL(url string, dbi *gorm.DB, forUpdate bool) (error, *URLContent) {
	if url == "" {
		return errors.New("url is empty"), nil
	}

	hashKey := GetURLHashKey(url)

	var urlContent URLContent

	sql := "SELECT * FROM url_contents WHERE hash_key = ?"

	if forUpdate && db.GetDbType() != db.SQLITE {
		sql = sql + " FOR UPDATE"
	}

	dbi.Raw(sql, hashKey).Scan(&urlContent)

	if urlContent.ID == 0 {
		return nil, nil
	}

	return nil, &urlContent
}
