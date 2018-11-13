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
	"encoding/base64"
	"errors"
	"math/big"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/primasio/wormhole/util"
	"golang.org/x/crypto/sha3"
)

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

func (user *User) VerifyPassword(password string) bool {
	return user.Password == user.getPasswordHash(password)
}

func (user *User) hashPassword() {
	user.Password = user.getPasswordHash(user.Password)
}

func (user *User) getPasswordHash(password string) string {

	// Password hash = sha3(sha3(password) + salt)

	hash := sha3.New256()
	hashByte := hash.Sum([]byte(password))

	saltByte := []byte(user.Salt)

	finalByteStr := append(hashByte, saltByte...)

	hash.Reset()
	finalByte := hash.Sum(finalByteStr)

	return base64.StdEncoding.EncodeToString(finalByte)
}

func (user *User) setSalt() {
	user.Salt = util.RandString(8)
}

func (user *User) BeforeCreate() error {

	user.CreatedAt = uint(time.Now().Unix())

	if user.Password != "" {
		user.setSalt()
		user.hashPassword()
	}

	user.SetBalance(big.NewInt(0))

	return nil
}

func (user *User) GetBalance() *big.Int {
	balanceNum := big.NewInt(0)
	balanceNum.SetString(user.Balance, 10)
	return balanceNum
}

func (user *User) SetBalance(num *big.Int) {
	user.Balance = num.String()
}

func (user *User) SetUniqueID(db *gorm.DB) error {
	var counter = 0

	for {
		counter = counter + 1
		uid := util.RandStringUppercase(8)

		check := &User{UniqueID: uid}

		db.Where(&check).First(&check)

		if check.ID == 0 {
			user.UniqueID = uid
			return nil
		}

		if counter >= 5 {
			// This is unlikely to happen
			// Must be error from other parts
			return errors.New("too many iterations while generating new session key")
		}
	}
}
