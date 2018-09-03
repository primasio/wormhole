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
	"github.com/primasio/wormhole/util"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func PrepareURLContentComment(content *models.URLContent) (error, *models.URLContentComment) {
	if systemUser.ID == 0 {
		return errors.New("system user not created"), nil
	}

	dbi := db.GetDb()

	randStr := util.RandString(10)

	urlContentComment := &models.URLContentComment{
		UserId:       systemUser.ID,
		Content:      "Comment " + randStr,
		URLContentId: content.ID,
	}

	urlContentComment.SetUniqueID(dbi)

	dbi.Create(&urlContentComment)

	return nil, urlContentComment
}

func TestURLContentCommentController_Create(t *testing.T) {
	PrepareAuthToken(t)

	err, urlContent := PrepareURLContent()
	assert.Equal(t, err, nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Set("url", urlContent.URL)
	form.Set("content", "<p>The comment of a URL.</p>")

	req, _ := http.NewRequest("POST", "/v1/urls/comments", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Authorization", authToken)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)
}

func TestURLContentCommentController_List(t *testing.T) {

	ResetDB()
	PrepareSystemUser()

	err, urlContent := PrepareURLContent()
	assert.Equal(t, err, nil)

	for i := 0; i < 30; i++ {
		err, _ := PrepareURLContentComment(urlContent)
		assert.Equal(t, err, nil)
	}

	url := url.QueryEscape(urlContent.URL)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/urls/comments?url="+url, nil)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)
}
