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
	"errors"
	"github.com/magiconair/properties/assert"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/models"
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

	urlContent := &models.URLContent{
		URL:      "https://cn.primas.io/12345",
		Title:    "Title of the Content",
		Content:  "<p>The content of the url.</p>",
		Abstract: "The content of the url.",
		UserId:   systemUser.ID,
	}

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
