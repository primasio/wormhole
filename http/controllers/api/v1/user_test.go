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
	"github.com/magiconair/properties/assert"
	"github.com/primasio/wormhole/http/server"
	"github.com/primasio/wormhole/tests"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestUserController_Create(t *testing.T) {
	tests.InitTestEnv("../../../../config/")
	router := server.NewRouter()
	w := httptest.NewRecorder()

	user, err := tests.CreateTestUser()

	if err != nil {
		log.Println(err)
		return
	}

	data := url.Values{}
	data.Set("Username", user.Username)
	data.Set("Password", user.Password)
	data.Set("Nickname", user.Nickname)

	req, _ := http.NewRequest("POST", "/v1/users", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)
}

func TestUserController_Auth(t *testing.T) {

}

func TestUserController_Get(t *testing.T) {

}
