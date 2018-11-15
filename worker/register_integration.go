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

package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/models"
	"github.com/primasio/wormhole/service"
)

// RegisterIntegrationWorker should only one instance in the world
type RegisterIntegrationWorker struct{}

func NewRegisterIntegrationWorker() *RegisterIntegrationWorker {
	return &RegisterIntegrationWorker{}
}

func (w *RegisterIntegrationWorker) Run() {
	if err := w.initLatestAssignedUserID(); err != nil {
		log.Fatal(err)
	}

	for {
		latestDoneUserID, err := w.getLatestAssignedUserID()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		latestRegisterUserID, err := w.getLatestRegisterUserID()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if latestDoneUserID < latestRegisterUserID {
			err := w.doAssignIntegration(latestDoneUserID, latestRegisterUserID)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func (w *RegisterIntegrationWorker) getLatestRegisterUserID() (uint, error) {
	dbi := db.GetDb()
	user := &models.User{}
	err := dbi.Last(user).Error
	return user.ID, err
}

func (w *RegisterIntegrationWorker) initLatestAssignedUserID() error {
	dbi := db.GetDb()
	info := &models.RegisterIntegrationWorkerInfo{}

	if err := dbi.Last(info).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			info.LastDoneUserID = 0
			if err := dbi.Create(info).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (w *RegisterIntegrationWorker) getLatestAssignedUserID() (uint, error) {
	dbi := db.GetDb()
	info := &models.RegisterIntegrationWorkerInfo{}

	if err := dbi.Last(info).Error; err != nil {
		return 0, err
	}

	return info.LastDoneUserID, nil
}

func (w *RegisterIntegrationWorker) doAssignIntegration(start, end uint) error {
	userID, err := w.getAssignUserID(start, end)
	if err != nil {
		return err
	}

	return w.assignUserIntegration(userID)
}

func (w *RegisterIntegrationWorker) getAssignUserID(start, end uint) (uint, error) {
	dbi := db.GetDb()
	user := &models.User{}

	if err := dbi.Where("id > ? AND id <= ?", start, end).First(user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (w *RegisterIntegrationWorker) assignUserIntegration(userID uint) error {
	// assign integration
	score := service.GetIntegration().GetRegisterScore()
	dbi := db.GetDb()
	tx := dbi.Begin()

	user := &models.User{}
	if err := db.ForUpdate(tx).Where("id = ?", userID).First(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	user.IncrementIntegration(score)
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// update latest done userid
	info := &models.RegisterIntegrationWorkerInfo{}
	if err := tx.Last(info).Error; err != nil {
		tx.Rollback()
		return err
	}

	info.LastDoneUserID = user.ID
	if err := tx.Save(info).Error; err != nil {
		tx.Rollback()
		return err
	}

	// insert integration history
	integrationHistory := &models.IntegrationHistory{UserID: user.ID, Integration: score}
	integrationHistory.Description = w.genIntegrationDescription(score)
	integrationHistory.Data = w.genIntegrationData(user.ID)
	integrationHistory.SetUniqueID()

	if err := tx.Create(integrationHistory).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (w *RegisterIntegrationWorker) genIntegrationDescription(score int64) string {
	return fmt.Sprintf("註冊獎勵積分: %d", score)
}

func (w *RegisterIntegrationWorker) genIntegrationData(userID uint) string {
	event := "USER_REGISTER"
	return fmt.Sprintf(`{"event": "%s", "user_id": %d}`, event, userID)
}
