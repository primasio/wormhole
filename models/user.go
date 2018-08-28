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

type User struct {
	BaseModel

	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password"`
	Salt     string `binding:"-"`
	Nickname string `form:"nickname" json:"nickname"`
}

func NewUser(username, password, nickname string) *User {

	user := &User{}

	user.Username = username
	user.Nickname = nickname

	user.Password = password

	return user
}

func (user *User) SetPassword(password string) {

}

func (user *User) setSalt() {

}
