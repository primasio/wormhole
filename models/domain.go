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
	"net/url"
)

type Domain struct {
	BaseModel
	UserId  uint   `json:"-"`
	Domain  string `gorm:"type:text" json:"domain"`
	Title   string `gorm:"type:text" json:"title"`
	HashKey string `gorm:"type:varchar(128);unique_index" json:"-"`

	IsActive bool `gorm:"default:false" json:"is_active"`
	Votes    uint `gorm:"default:1" json:"votes"`
}

func ExtractDomainFromURL(urlStr string) (error, string) {
	u, err := url.Parse(urlStr)

	if err != nil {
		return err, ""
	}

	return nil, u.Host
}

func CleanDomain(domain string) string {
	//TODO: Remove trailing slash, remove www, etc
	return domain
}

func GetDomainByDomainName(domain string, dbi *gorm.DB, forUpdate bool) (error, *Domain) {

	if domain == "" {
		return errors.New("domain is empty"), nil
	}

	hashKey := GetDomainHashKey(domain)

	var domainModel Domain

	sql := "SELECT * FROM domains WHERE hash_key = ?"

	if forUpdate && db.GetDbType() != db.SQLITE {
		sql = sql + " FOR UPDATE"
	}

	dbi.Raw(sql, hashKey).Scan(&domainModel)

	if domainModel.ID == 0 {
		return nil, nil
	}

	return nil, &domainModel
}

func GetDomainHashKey(domain string) string {

	sumBytes := sha1.Sum([]byte(domain))

	return hex.EncodeToString(sumBytes[:])
}
