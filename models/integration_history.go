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
	"errors"
	"fmt"
	"io"
)

type IntegrationHistory struct {
	BaseModel

	UserID      uint   `json:"-"`
	UniqueID    string `json:"id" gorm:"type:varchar(128);unique_index"`
	Integration int64  `json:"integration"`
	Description string `json:"description"`
	Data        string `json:"-"`

	User User `gorm:"save_associations:false" json:"user"`
}

func (m *IntegrationHistory) SetUniqueID() error {
	if m.Data == "" {
		return errors.New("Integration History Data Required")
	}

	h := sha1.New()
	io.WriteString(h, m.Data)
	m.UniqueID = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}
