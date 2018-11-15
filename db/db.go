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

package db

import (
	"io/ioutil"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/primasio/wormhole/config"
)

const (
	SQLITE = "sqlite3"
	MYSQL  = "mysql"
)

var instance *gorm.DB
var instanceType string

func GetDb() *gorm.DB {
	return instance
}

func GetDbType() string {
	return instanceType
}

func Init() error {

	c := config.GetConfig()

	instanceType = c.GetString("db.type")
	dbConn := c.GetString("db.connection")

	if dbConn == "" {
		f, err := ioutil.TempFile("", "")
		if err != nil {
			return err
		}
		dbConn := f.Name()
		f.Close()
		os.Remove(dbConn)
	}

	var err error

	instance, err = gorm.Open(instanceType, dbConn)

	if err != nil {
		return err
	}

	instance.Set("gorm:table_options", "charset=utf8mb4")

	return nil
}

func ForUpdate(tx *gorm.DB) *gorm.DB {
	if GetDbType() != SQLITE {
		return tx.Set("gorm:query_option", "FOR UPDATE")
	}
	return tx
}
