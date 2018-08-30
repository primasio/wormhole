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

package v1_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/http/server"
	"github.com/primasio/wormhole/models"
	"github.com/primasio/wormhole/tests"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
)

var router *gin.Engine
var systemUser *models.User
var authToken string

func TestMain(m *testing.M) {
	before()
	retCode := m.Run()
	os.Exit(retCode)
}

func PrepareSystemUser() {
	if systemUser != nil {
		return
	}

	user, err := tests.CreateTestUser()

	if err != nil {
		log.Fatal(err)
	}

	dbi := db.GetDb()

	user.SetUniqueID(dbi)

	dbi.Create(&user)

	systemUser = user
}

func PrepareAuthToken(t *testing.T) {

	PrepareSystemUser()

	if authToken != "" {
		return
	}

	w2 := httptest.NewRecorder()

	login := url.Values{}
	login.Set("username", systemUser.Username)
	login.Set("password", "PrimasGoGoGo")
	login.Set("remember", "on")

	req2, _ := http.NewRequest("POST", "/v1/users/auth", strings.NewReader(login.Encode()))
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Add("Content-Length", strconv.Itoa(len(login.Encode())))

	router.ServeHTTP(w2, req2)

	responseStr := w2.Body.String()

	log.Println(responseStr)
	assert.Equal(t, w2.Code, 200)

	var returnData map[string]*json.RawMessage

	err := json.Unmarshal([]byte(responseStr), &returnData)
	assert.Equal(t, err, nil)

	var tokenStruct map[string]string

	err = json.Unmarshal(*returnData["data"], &tokenStruct)
	assert.Equal(t, err, nil)

	authToken = tokenStruct["token"]

	log.Println("token: " + authToken)
}

func before() {
	log.Println("Setting up test environment")
	tests.InitTestEnv("../../../../config/")
	router = server.NewRouter()
}
