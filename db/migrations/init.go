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

package migrations

import (
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"github.com/primasio/wormhole/db"
	"gopkg.in/gormigrate.v1"
)

func Migrate() error {

	mgs := getMigrations()

	dbi := db.GetDb()

	m := gormigrate.New(dbi, gormigrate.DefaultOptions, mgs)

	if err := m.Migrate(); err != nil {
		glog.Info("Migration failed")
		return err
	}

	glog.Info("Migration success")

	return nil
}

func getMigrations() []*gormigrate.Migration {

	migrations := initialTables()

	migrations = append(migrations, Migration20180911()...)

	return migrations
}

func initialTables() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "201809111325",
			Migrate: func(tx *gorm.DB) error {

				// it's a good practice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time

				type BaseModel struct {
					ID        uint `gorm:"primary_key" json:"-"`
					CreatedAt uint `json:"created_at"`
					UpdatedAt uint `json:"updated_at"`
				}

				type User struct {
					BaseModel
					UniqueID string `json:"id" gorm:"type:varchar(128);unique_index"`
					Username string `json:"-" gorm:"type:varchar(128);index"`
					Password string `json:"-"`
					Salt     string `json:"-"`
					Nickname string `json:"nickname"`
					Balance  string `json:"balance"`
				}

				type UserOAuth struct {
					BaseModel
					UserID     uint   `gorm:"index"`
					VendorType uint   `gorm:"index"`
					VendorID   string `gorm:"index"`
				}

				if err := tx.AutoMigrate(&User{}).Error; err != nil {
					return err
				}

				if err := tx.AutoMigrate(&UserOAuth{}).Error; err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.DropTable("users").Error; err != nil {
					return err
				}

				if err := tx.DropTable("user_o_auths").Error; err != nil {
					return err
				}

				return nil
			},
		},
	}
}
