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
	"errors"
	"github.com/magiconair/properties/assert"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/http/token"
	"github.com/primasio/wormhole/models"
	"github.com/primasio/wormhole/tests"
	"github.com/primasio/wormhole/util"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func PrepareURLContent() (error, *models.URLContent) {

	if systemUser.ID == 0 {
		return errors.New("system user not created"), nil
	}

	randStr := util.RandString(10)

	urlContent := &models.URLContent{
		URL:      "https://cn.primas.io/12345" + randStr,
		Title:    "Title of the Content",
		Content:  "<p>The content of the url.</p>",
		Abstract: "The content of the url.",
		UserId:   systemUser.ID,
	}

	urlContent.HashKey = models.GetURLHashKey(urlContent.URL)

	dbi := db.GetDb()
	dbi.Create(&urlContent)

	return nil, urlContent
}

func TestURLContentController_Create(t *testing.T) {
	PrepareAuthToken(t)

	data := url.Values{}
	data.Set("url", "https://cn.primas.io/abcdefg")
	data.Set("title", "Title of the Content")
	data.Set("content", "<p>The content of the url.</p>")
	data.Set("abstract", "The content of the url.")

	req, _ := http.NewRequest("POST", "/v1/urls", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Authorization", authToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)

	// Create again to see error

	req2, _ := http.NewRequest("POST", "/v1/urls", strings.NewReader(data.Encode()))
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req2.Header.Add("Authorization", authToken)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	log.Println(w2.Body.String())
	assert.Equal(t, w2.Code, 400)
}

func TestURLContentController_Get(t *testing.T) {
	PrepareAuthToken(t)
	err, urlContent := PrepareURLContent()

	assert.Equal(t, err, nil)

	escaped := url.QueryEscape(urlContent.URL)

	log.Println("escaped url: " + escaped)

	req, _ := http.NewRequest("GET", "/v1/urls/url?url="+escaped, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)
}

func TestURLContentController_Vote(t *testing.T) {
	PrepareAuthToken(t)
	err, urlContent := PrepareURLContent()

	assert.Equal(t, err, nil)

	escaped := url.QueryEscape(urlContent.URL)

	log.Println("escaped url: " + escaped)

	// Author cannot vote

	req, _ := http.NewRequest("PUT", "/v1/urls/url?url="+escaped, nil)
	req.Header.Add("Authorization", authToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 400)

	dbi := db.GetDb()

	// Vote using another user
	user, err := tests.CreateTestUser()
	assert.Equal(t, err, nil)

	err = user.SetUniqueID(dbi)
	assert.Equal(t, err, nil)
	dbi.Create(&user)

	err, userToken := token.IssueToken(user.ID, false)
	assert.Equal(t, err, nil)

	req2, _ := http.NewRequest("PUT", "/v1/urls/url?url="+escaped, nil)
	req2.Header.Add("Authorization", userToken.Token)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	log.Println(w2.Body.String())
	assert.Equal(t, w2.Code, 200)

	// Vote again
	w3 := httptest.NewRecorder()

	req3, _ := http.NewRequest("PUT", "/v1/urls/url?url="+escaped, nil)
	req3.Header.Add("Authorization", userToken.Token)

	router.ServeHTTP(w3, req3)

	log.Println(w3.Body.String())
	assert.Equal(t, w3.Code, 400)

	// Test URL that is already active

	urlContent.IsActive = true

	dbi.Save(&urlContent)

	req4, _ := http.NewRequest("PUT", "/v1/urls/url?url="+escaped, nil)
	req4.Header.Add("Authorization", authToken)

	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	log.Println(w4.Body.String())
	assert.Equal(t, w4.Code, 400)
}

func TestURLContentController_List(t *testing.T) {
	ResetDB()
	PrepareSystemUser()

	// Create a list of url content

	urlContents := make([]*models.URLContent, 10)

	dbi := db.GetDb()

	largestActiveCreatedAt := uint(0)
	largestVotingCreatedAt := uint(0)

	for i := 0; i < 10; i++ {
		err, urlContent := PrepareURLContent()
		assert.Equal(t, err, nil)

		if i%2 == 0 {
			urlContent.IsActive = true
			largestActiveCreatedAt++
			urlContent.CreatedAt = largestActiveCreatedAt
		} else {
			largestVotingCreatedAt++
			urlContent.CreatedAt = largestVotingCreatedAt
		}

		dbi.Save(&urlContent)
		urlContents[i] = urlContent
	}

	// Get voting list
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/v1/urls?type=voting", nil)

	router.ServeHTTP(w, req)

	resp := w.Body.String()

	log.Println(resp)
	assert.Equal(t, w.Code, 200)

	// Check result
	urlContentList := getUrlContentListFromJsonString(resp, t)
	assert.Equal(t, len(urlContentList), 5)
	assert.Equal(t, urlContentList[0].CreatedAt, largestVotingCreatedAt)
	assert.Equal(t, urlContentList[0].IsActive, false)

	// Get active list
	w2 := httptest.NewRecorder()

	req2, _ := http.NewRequest("GET", "/v1/urls?type=active", nil)

	router.ServeHTTP(w2, req2)

	resp2 := w2.Body.String()

	log.Println(resp2)
	assert.Equal(t, w2.Code, 200)

	// Check result
	urlContentList2 := getUrlContentListFromJsonString(resp2, t)
	assert.Equal(t, len(urlContentList2), 5)
	assert.Equal(t, urlContentList2[0].CreatedAt, largestActiveCreatedAt)
	assert.Equal(t, urlContentList2[0].IsActive, true)
}

func getUrlContentListFromJsonString(jsonStr string, t *testing.T) []models.URLContent {

	var returnData map[string]*json.RawMessage

	err := json.Unmarshal([]byte(jsonStr), &returnData)
	assert.Equal(t, err, nil)

	var list []models.URLContent

	err = json.Unmarshal(*returnData["data"], &list)
	assert.Equal(t, err, nil)

	return list
}
