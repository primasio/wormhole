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
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestURLContentCommentVoteController_Create(t *testing.T) {
	PrepareAuthToken(t)

	user, err := PrepareTestUser()
	assert.Equal(t, err, nil)

	urlContent, err := PrepareURLContentWithUser(user)
	assert.Equal(t, err, nil)

	urlContentComment, err := PrepareURLContentCommentWithContent(urlContent)
	assert.Equal(t, err, nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Set("like", "true")

	req, _ := http.NewRequest("POST", "/v1/comments/"+urlContentComment.UniqueID+"/votes", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Authorization", authToken)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)
}

func TestURLContentCommentVoteController_Update(t *testing.T) {
	PrepareAuthToken(t)

	user, err := PrepareTestUser()
	assert.Equal(t, err, nil)

	urlContent, err := PrepareURLContentWithUser(user)
	assert.Equal(t, err, nil)

	urlContentComment, err := PrepareURLContentCommentWithContent(urlContent)
	assert.Equal(t, err, nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Set("like", "true")

	req, _ := http.NewRequest("POST", "/v1/comments/"+urlContentComment.UniqueID+"/votes", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Authorization", authToken)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)

	// Update
	form = url.Values{}
	form.Set("like", "false")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/v1/comments/"+urlContentComment.UniqueID+"/votes", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Authorization", authToken)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)
}

func TestURLContentCommentVoteController_Delete(t *testing.T) {
	PrepareAuthToken(t)

	user, err := PrepareTestUser()
	assert.Equal(t, err, nil)

	urlContent, err := PrepareURLContentWithUser(user)
	assert.Equal(t, err, nil)

	urlContentComment, err := PrepareURLContentCommentWithContent(urlContent)
	assert.Equal(t, err, nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Set("like", "true")

	req, _ := http.NewRequest("POST", "/v1/comments/"+urlContentComment.UniqueID+"/votes", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Authorization", authToken)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)

	// Delete
	form = url.Values{}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/v1/comments/"+urlContentComment.UniqueID+"/votes", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Authorization", authToken)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)

	// Update
	form = url.Values{}
	form.Set("like", "false")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/v1/comments/"+urlContentComment.UniqueID+"/votes", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Authorization", authToken)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 400)
}
