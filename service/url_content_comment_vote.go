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

package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/models"
)

var uccVote *URLContentCommentVote
var uccVoteOnce sync.Once

type URLContentCommentVote struct{}

func GetURLContentCommentVote() *URLContentCommentVote {
	uccVoteOnce.Do(func() {
		uccVote = &URLContentCommentVote{}
	})

	return uccVote
}

func (s *URLContentCommentVote) CreateVote(dbi *gorm.DB, comment *models.URLContentComment, user *models.User, like bool) error {

	vote := &models.URLContentCommentVote{UserID: user.ID, URLContentCommentID: comment.ID, Like: like}
	vote.SetUniqueID()

	tx := dbi.Begin()
	if err := tx.Create(vote).Error; err != nil {
		tx.Rollback()
		return err
	}

	contentComment := &models.URLContentComment{}
	if err := db.ForUpdate(tx).Where("id = ?", comment.ID).First(contentComment).Error; err != nil {
		tx.Rollback()
		return err
	}

	contentComment.IncrementVote(like)
	if err := tx.Save(contentComment).Error; err != nil {
		tx.Rollback()
		return err
	}

	commentOwner := &models.User{}
	if err := db.ForUpdate(tx).Where("id = ?", comment.UserID).First(commentOwner).Error; err != nil {
		tx.Rollback()
		return err
	}

	score := GetIntegration().GetURLContentCommentVoteScore(like)

	commentOwner.IncrementCommentVote(like)
	commentOwner.IncrementIntegration(score)

	if err := tx.Save(commentOwner).Error; err != nil {
		tx.Rollback()
		return err
	}

	integrationHistory := &models.IntegrationHistory{UserID: commentOwner.ID, Integration: score}
	integrationHistory.Description = s.GenIntegrationDescription(user.Nickname, score, like)
	integrationHistory.Data = s.GenIntegrationData(user.ID, comment.ID, vote.ID)
	integrationHistory.SetUniqueID()

	if err := tx.Create(integrationHistory).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *URLContentCommentVote) UpdateVote(dbi *gorm.DB, comment *models.URLContentComment, user *models.User, like bool) error {
	vote := &models.URLContentCommentVote{UserID: user.ID, URLContentCommentID: comment.ID, Like: like}
	vote.SetUniqueID()

	exists, err := vote.CheckVoteExists(dbi, vote.UniqueID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("Vote not found")
	}

	tx := dbi.Begin()
	oldVote := &models.URLContentCommentVote{}
	err = db.ForUpdate(tx).Where("unique_id = ?", vote.UniqueID).First(oldVote).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if oldVote.Like == vote.Like {
		tx.Rollback()
		return nil
	}

	oldVote.Like = vote.Like
	if err := tx.Save(oldVote).Error; err != nil {
		tx.Rollback()
		return err
	}

	contentComment := &models.URLContentComment{}
	if err := db.ForUpdate(tx).Where("id = ?", comment.ID).First(contentComment).Error; err != nil {
		tx.Rollback()
		return err
	}

	contentComment.SwitchVote(like)
	if err := tx.Save(contentComment).Error; err != nil {
		tx.Rollback()
		return err
	}

	commentOwner := &models.User{}
	if err := db.ForUpdate(tx).Where("id = ?", comment.UserID).First(commentOwner).Error; err != nil {
		tx.Rollback()
		return err
	}

	oldScore := GetIntegration().GetURLContentCommentVoteScore(!like)
	score := GetIntegration().GetURLContentCommentVoteScore(like)

	commentOwner.SwitchCommentVote(like)
	commentOwner.IncrementIntegration(-oldScore)
	commentOwner.IncrementIntegration(score)

	if err := tx.Save(commentOwner).Error; err != nil {
		tx.Rollback()
		return err
	}

	data := s.GenIntegrationData(user.ID, comment.ID, oldVote.ID)

	h := sha1.New()
	io.WriteString(h, data)
	uid := fmt.Sprintf("%x", h.Sum(nil))

	integrationHistory := &models.IntegrationHistory{}
	if err := db.ForUpdate(tx).Where("unique_id = ?", uid).First(integrationHistory).Error; err != nil {
		tx.Rollback()
		return err
	}

	integrationHistory.Description = s.GenIntegrationDescription(user.Nickname, score, like)
	integrationHistory.Integration = score

	if err := tx.Save(integrationHistory).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

func (s *URLContentCommentVote) CancelVote(dbi *gorm.DB, comment *models.URLContentComment, user *models.User) error {
	vote := &models.URLContentCommentVote{UserID: user.ID, URLContentCommentID: comment.ID}
	vote.SetUniqueID()

	exists, err := vote.CheckVoteExists(dbi, vote.UniqueID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("Vote not found")
	}

	tx := dbi.Begin()
	oldVote := &models.URLContentCommentVote{}
	err = tx.Where("unique_id = ?", vote.UniqueID).First(oldVote).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(oldVote).Error; err != nil {
		tx.Rollback()
		return err
	}

	contentComment := &models.URLContentComment{}
	if err := db.ForUpdate(tx).Where("id = ?", comment.ID).First(contentComment).Error; err != nil {
		tx.Rollback()
		return err
	}

	contentComment.CancelVote(oldVote.Like)
	if err := tx.Save(contentComment).Error; err != nil {
		tx.Rollback()
		return err
	}

	commentOwner := &models.User{}
	if err := db.ForUpdate(tx).Where("id = ?", comment.UserID).First(commentOwner).Error; err != nil {
		tx.Rollback()
		return err
	}

	score := GetIntegration().GetURLContentCommentVoteScore(oldVote.Like)

	commentOwner.CancelCommentVote(oldVote.Like)
	commentOwner.IncrementIntegration(-score)

	if err := tx.Save(commentOwner).Error; err != nil {
		tx.Rollback()
		return err
	}

	data := s.GenIntegrationData(user.ID, comment.ID, oldVote.ID)

	h := sha1.New()
	io.WriteString(h, data)
	uid := fmt.Sprintf("%x", h.Sum(nil))

	integrationHistory := &models.IntegrationHistory{}
	if err := tx.Where("unique_id = ?", uid).First(integrationHistory).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(integrationHistory).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *URLContentCommentVote) GenIntegrationDescription(nickname string, score int64, like bool) string {
	if like {
		return fmt.Sprintf(`%s 爲你點讚, 獎勵積分 %d`, nickname, score)
	}
	return fmt.Sprintf(`%s 鄙視了你, %d 積分受到傷害`, nickname, score)
}

func (s *URLContentCommentVote) GenIntegrationData(userID, urlContentCommentID, urlContentCommentVoteID uint) string {
	event := "URL_CONTENT_COMMENT_VOTE"
	return fmt.Sprintf(`{"event": "%s", "user_id": %d, "url_content_comment_id": %d, "url_content_comment_vote_id": %d}`, event, userID, urlContentCommentID, urlContentCommentVoteID)
}
