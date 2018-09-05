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

package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/primasio/wormhole/cache"
	"github.com/primasio/wormhole/config"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/http/server"
	"github.com/primasio/wormhole/models"
	"os"
)

func main() {

	flag.Parse()

	// Init Environment
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init Config
	configDir := os.Getenv("APP_CONFIG")

	var err error

	if configDir != "" {
		err = config.Init(env, &configDir)
	} else {
		err = config.Init(env, nil)
	}

	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}

	// Init Database
	if err := db.Init(); err != nil {
		glog.Error(err)
		os.Exit(1)
	}

	if env == "development" {
		models.AutoMigrateModels()
	}

	// Init Cache
	if err := cache.InitCache(); err != nil {
		glog.Error(err)
		os.Exit(1)
	}

	// Start HTTP server
	server.Init()
}
