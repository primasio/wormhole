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

	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

type URLContentCommentVote struct {
	BaseModel
	UniqueID            string `json:"id" gorm:"type:varchar(128);unique_index"`
	UserID              uint   `json:"-"`
	URLContentCommentID uint   `json:"-"`
	Like                bool   `json:"like"`

	User User `gorm:"save_associations:false" json:"user"`
}

func (comment *URLContentCommentVote) SetUniqueID() error {

	if comment.UserID == 0 || comment.URLContentCommentID == 0 {
		return errors.New("UserID Or URLContentCommentID Zero")
	}

	h := sha1.New()
	io.WriteString(h, fmt.Sprintf("%d%d", comment.UserID, comment.URLContentCommentID))
	comment.UniqueID = fmt.Sprintf("%x", h.Sum(nil))

	return nil

}

func (common *URLContentCommentVote) CheckVoteExists(dbi *gorm.DB, uniqueID string) (bool, error) {
	vote := &URLContentCommentVote{}
	err := dbi.Where("unique_id = ?", uniqueID).First(vote).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	logrus.Errorf("%v+", vote)

	return vote.ID != 0, nil
}
