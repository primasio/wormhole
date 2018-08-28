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
	"encoding/hex"
	"github.com/primasio/wormhole/util"
	"golang.org/x/crypto/sha3"
	"time"
)

type User struct {
	BaseModel
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Salt     string `binding:"-"`
	Nickname string `form:"nickname" json:"nickname" binding:"required"`
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

	return hex.EncodeToString(finalByte)
}

func (user *User) setSalt() {
	user.Salt = util.RandString(8)
}

func (user *User) BeforeCreate() error {

	user.CreatedAt = uint(time.Now().Unix())

	user.setSalt()
	user.hashPassword()

	return nil
}