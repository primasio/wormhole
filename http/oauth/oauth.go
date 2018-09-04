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

package oauth

import (
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/models"
	"log"
)

type OAuthResult struct {
	Type      uint
	Id        string
	Email     string
	Name      string
	AvatarURL string
}

func (oauthResult *OAuthResult) Process() (err error, userId uint) {

	// Find the corresponding user in our DB
	// or create one if not exists

	userOAuth := &models.UserOAuth{
		VendorType: oauthResult.Type,
		VendorID:   oauthResult.Id,
	}

	dbi := db.GetDb()
	dbi.Where(&userOAuth).First(&userOAuth)

	if userOAuth.ID != 0 {
		// User exists
		return nil, userOAuth.UserID
	}

	// User not exists, create a new one

	user := &models.User{}
	user.Username = ""
	user.Password = ""

	log.Println("OAuth name: " + oauthResult.Name)
	log.Println("OAuth email: " + oauthResult.Email)

	if oauthResult.Name != "" {
		user.Nickname = oauthResult.Name
	} else {
		user.Nickname = oauthResult.Email
	}

	if err := user.SetUniqueID(dbi); err != nil {
		return err, 0
	}

	// Use transaction to ensure consistency for user & oauth credential creation

	tx := dbi.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err, 0
	}

	if err := tx.Create(&userOAuth).Error; err != nil {
		tx.Rollback()
		return err, 0
	}

	tx.Commit()

	if oauthResult.AvatarURL != "" {
		// TODO: Async downloading avatar if exists
	}

	return nil, user.ID
}
