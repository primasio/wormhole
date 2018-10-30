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
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/models"
	"github.com/primasio/wormhole/util"
)

func PrepareURLContent() (error, *models.URLContent) {

	PrepareSystemUser()

	err, domain := PrepareDomain()

	if err != nil {
		return err, nil
	}

	randStr := util.RandString(10)

	urlContent := &models.URLContent{
		URL:    "https://" + domain.Domain + "/12345" + randStr,
		UserID: systemUser.ID,
	}

	urlContent.HashKey = models.GetURLHashKey(urlContent.URL)

	dbi := db.GetDb()
	dbi.Create(&urlContent)

	return nil, urlContent
}

func TestURLContentController_Get(t *testing.T) {

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
