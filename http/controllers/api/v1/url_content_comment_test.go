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
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/models"
	"github.com/primasio/wormhole/util"
)

func PrepareURLContentComment(content *models.URLContent) (error, *models.URLContentComment) {
	if systemUser.ID == 0 {
		return errors.New("system user not created"), nil
	}

	dbi := db.GetDb()

	randStr := util.RandString(10)

	urlContentComment := &models.URLContentComment{
		UserID:       systemUser.ID,
		Content:      "Comment " + randStr,
		URLContentId: content.ID,
	}

	err := urlContentComment.SetUniqueID(dbi)

	if err != nil {
		return err, nil
	}

	dbi.Create(&urlContentComment)

	content.TotalComment++

	dbi.Save(content)

	return nil, urlContentComment
}

func PrepareURLContentCommentWithContent(content *models.URLContent) (*models.URLContentComment, error) {
	dbi := db.GetDb()

	randStr := util.RandString(10)

	urlContentComment := &models.URLContentComment{
		UserID:       content.UserID,
		Content:      "Comment " + randStr,
		URLContentId: content.ID,
	}

	err := urlContentComment.SetUniqueID(dbi)

	if err != nil {
		return nil, err
	}

	dbi.Create(&urlContentComment)

	content.TotalComment++

	dbi.Save(content)

	return urlContentComment, nil
}

func TestURLContentCommentController_Create(t *testing.T) {
	PrepareAuthToken(t)

	err, urlContent := PrepareURLContent()
	assert.Equal(t, err, nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Set("url", urlContent.URL)
	form.Set("content", "<p>The comment of a URL.</p>")

	req, _ := http.NewRequest("POST", "/v1/comments", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Authorization", authToken)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)
}

func TestURLContentCommentController_List(t *testing.T) {

	ResetDB()

	err, urlContent := PrepareURLContent()
	assert.Equal(t, err, nil)

	for i := 0; i < 30; i++ {
		err, _ := PrepareURLContentComment(urlContent)
		assert.Equal(t, err, nil)
	}

	urlEscaped := url.QueryEscape(urlContent.URL)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/comments?url="+urlEscaped, nil)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)

	// Test list empty url

	urlNull := url.QueryEscape("http://not.exist.com/a/web/page.html")

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/v1/comments?url="+urlNull, nil)

	router.ServeHTTP(w2, req2)

	log.Println(w2.Body.String())
	assert.Equal(t, w2.Code, 200)
}

func TestURLContentCommentController_ListWithVote(t *testing.T) {

	ResetDB()

	err, urlContent := PrepareURLContent()
	assert.Equal(t, err, nil)

	for i := 0; i < 30; i++ {
		err, _ := PrepareURLContentComment(urlContent)
		assert.Equal(t, err, nil)
	}

	urlEscaped := url.QueryEscape(urlContent.URL)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/authorized/comments?url="+urlEscaped, nil)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 401)

	PrepareAuthToken(t)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/authorized/comments?url="+urlEscaped, nil)
	req.Header.Add("Authorization", authToken)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)

}

func TestURLContentCommentController_Delete(t *testing.T) {

	PrepareAuthToken(t)

	err, urlContent := PrepareURLContent()
	assert.Equal(t, err, nil)

	comments := make([]*models.URLContentComment, 10)

	for i := 0; i < 10; i++ {
		err, comment := PrepareURLContentComment(urlContent)
		assert.Equal(t, err, nil)

		comments[i] = comment
	}

	// Delete first comment

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/comments/"+comments[0].UniqueID+"?token=TokenNotNeededForTest", nil)
	req.Header.Add("Authorization", authToken)

	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)

	// Delete second comment

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/v1/comments/"+comments[1].UniqueID, nil)
	req2.Header.Add("Authorization", authToken)

	router.ServeHTTP(w2, req2)

	log.Println(w.Body.String())
	assert.Equal(t, w2.Code, 200)
}
