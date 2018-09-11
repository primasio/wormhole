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
	"github.com/primasio/wormhole/config"
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

func PrepareDomain() (error, *models.Domain) {

	if systemUser.ID == 0 {
		return errors.New("system user not created"), nil
	}

	randStr := util.RandString(10)

	domainModel := &models.Domain{
		Domain: "primas" + randStr + ".io",
		UserID: systemUser.ID,
	}

	domainModel.HashKey = models.GetDomainHashKey(domainModel.Domain)
	domainModel.IsActive = true

	dbi := db.GetDb()
	dbi.Create(&domainModel)

	return nil, domainModel
}

func TestDomainController_Create(t *testing.T) {
	PrepareAuthToken(t)

	data := url.Values{}
	data.Set("domain", "cn.primas.io")
	data.Set("title", "Title of the Domain")

	req, _ := http.NewRequest("POST", "/v1/domains", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Authorization", authToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)

	// Create again to see error

	req2, _ := http.NewRequest("POST", "/v1/domains", strings.NewReader(data.Encode()))
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req2.Header.Add("Authorization", authToken)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	log.Println(w2.Body.String())
	assert.Equal(t, w2.Code, 400)
}

func TestDomainController_Get(t *testing.T) {
	PrepareAuthToken(t)
	err, domainModel := PrepareDomain()

	assert.Equal(t, err, nil)

	escaped := url.QueryEscape(domainModel.Domain)

	log.Println("escaped domain: " + escaped)

	req, _ := http.NewRequest("GET", "/v1/domains/domain?domain="+escaped, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	log.Println(w.Body.String())
	assert.Equal(t, w.Code, 200)
}

func TestDomainController_Vote(t *testing.T) {
	PrepareAuthToken(t)
	err, domainModel := PrepareDomain()

	domainModel.IsActive = false
	db.GetDb().Save(domainModel)

	assert.Equal(t, err, nil)

	escaped := url.QueryEscape(domainModel.Domain)

	log.Println("escaped domain: " + escaped)

	// Author cannot vote

	req, _ := http.NewRequest("PUT", "/v1/domains/domain?domain="+escaped, nil)
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

	req2, _ := http.NewRequest("PUT", "/v1/domains/domain?domain="+escaped, nil)
	req2.Header.Add("Authorization", userToken.Token)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	log.Println(w2.Body.String())
	assert.Equal(t, w2.Code, 200)

	// Vote again
	w3 := httptest.NewRecorder()

	req3, _ := http.NewRequest("PUT", "/v1/domains/domain?domain="+escaped, nil)
	req3.Header.Add("Authorization", userToken.Token)

	router.ServeHTTP(w3, req3)

	log.Println(w3.Body.String())
	assert.Equal(t, w3.Code, 400)

	// Test domain that is already active

	domainModel.IsActive = true

	dbi.Save(&domainModel)

	req4, _ := http.NewRequest("PUT", "/v1/domains/domain?domain="+escaped, nil)
	req4.Header.Add("Authorization", authToken)

	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	log.Println(w4.Body.String())
	assert.Equal(t, w4.Code, 400)
}

func TestDomainController_List(t *testing.T) {
	ResetDB()
	PrepareSystemUser()

	// Create a list of url content

	domains := make([]*models.Domain, 10)

	dbi := db.GetDb()

	largestActiveCreatedAt := uint(0)
	largestVotingCreatedAt := uint(0)

	for i := 0; i < 10; i++ {
		err, domain := PrepareDomain()
		assert.Equal(t, err, nil)

		domain.IsActive = false
		db.GetDb().Save(domain)

		if i%2 == 0 {
			domain.IsActive = true
			largestActiveCreatedAt++
			domain.CreatedAt = largestActiveCreatedAt
		} else {
			largestVotingCreatedAt++
			domain.CreatedAt = largestVotingCreatedAt
		}

		dbi.Save(&domain)
		domains[i] = domain
	}

	// Get voting list
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/v1/domains?type=voting", nil)

	router.ServeHTTP(w, req)

	resp := w.Body.String()

	log.Println(resp)
	assert.Equal(t, w.Code, 200)

	// Check result
	domainList := getDomainListFromJsonString(resp, t)
	assert.Equal(t, len(domainList), 5)
	assert.Equal(t, domainList[0].CreatedAt, largestVotingCreatedAt)
	assert.Equal(t, domainList[0].IsActive, false)

	// Get active list
	w2 := httptest.NewRecorder()

	req2, _ := http.NewRequest("GET", "/v1/domains?type=active", nil)

	router.ServeHTTP(w2, req2)

	resp2 := w2.Body.String()

	log.Println(resp2)
	assert.Equal(t, w2.Code, 200)

	// Check result
	domainList2 := getDomainListFromJsonString(resp2, t)
	assert.Equal(t, len(domainList2), 5)
	assert.Equal(t, domainList2[0].CreatedAt, largestActiveCreatedAt)
	assert.Equal(t, domainList2[0].IsActive, true)
}

func TestDomainController_Approve(t *testing.T) {
	PrepareAuthToken(t)

	err, domainModel := PrepareDomain()
	assert.Equal(t, err, nil)

	domainModel.IsActive = false
	db.GetDb().Save(domainModel)

	escaped := url.QueryEscape(domainModel.Domain)

	req, _ := http.NewRequest("POST", "/v1/domains/domain/approval?domain="+escaped, nil)
	req.Header.Add("Authorization", config.GetConfig().GetString("admin.key"))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Body.String()

	log.Println(resp)
	assert.Equal(t, w.Code, 200)
}

func getDomainListFromJsonString(jsonStr string, t *testing.T) []models.Domain {

	var returnData map[string]*json.RawMessage

	err := json.Unmarshal([]byte(jsonStr), &returnData)
	assert.Equal(t, err, nil)

	var list []models.Domain

	err = json.Unmarshal(*returnData["data"], &list)
	assert.Equal(t, err, nil)

	return list
}
