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

	// Test normal creation

	user, err := tests.CreateTestUser()

	if err != nil {
		log.Println(err)
		return
	}

	data := url.Values{}
	data.Set("username", user.Username)
	data.Set("password", user.Password)
	data.Set("nickname", user.Nickname)

	req, _ := http.NewRequest("POST", "/v1/users", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)

	// Test duplicated user creation

	w2 := httptest.NewRecorder()

	req2, _ := http.NewRequest("POST", "/v1/users", strings.NewReader(data.Encode()))
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	router.ServeHTTP(w2, req2)

	log.Println(w2.Body.String())
	assert.Equal(t, w2.Code, 400)
}

func TestUserController_Auth(t *testing.T) {

	tests.InitTestEnv("../../../../config/")
	router := server.NewRouter()
	w := httptest.NewRecorder()

	// Test normal creation

	user, err := tests.CreateTestUser()

	if err != nil {
		log.Println(err)
		return
	}

	data := url.Values{}
	data.Set("username", user.Username)
	data.Set("password", user.Password)
	data.Set("nickname", user.Nickname)

	req, _ := http.NewRequest("POST", "/v1/users", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)

	// Login using this user

	w2 := httptest.NewRecorder()

	login := url.Values{}
	login.Set("username", user.Username)
	login.Set("password", user.Password)
	login.Set("remember", "on")

	req2, _ := http.NewRequest("POST", "/v1/users/auth", strings.NewReader(login.Encode()))
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Add("Content-Length", strconv.Itoa(len(login.Encode())))

	router.ServeHTTP(w2, req2)

	log.Println(w2.Body.String())
	assert.Equal(t, w2.Code, 200)
}

func TestUserController_Get(t *testing.T) {

}
