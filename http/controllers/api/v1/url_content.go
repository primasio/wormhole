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

package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/models"
)

type URLContentController struct{}

func (ctrl *URLContentController) Get(c *gin.Context) {

	url := c.Query("url")

	if url == "" {
		Error("missing query param url", c)
		return
	}

	// Check whether the domain is approved.

	dbi := db.GetDb()

	cleanedUrl := models.CleanURL(url)

	err, domain := models.ExtractDomainFromURL(cleanedUrl)

	if err != nil {
		Error(err.Error(), c)
		return
	}

	err, domainModel := models.GetDomainByDomainName(domain, dbi, false)

	if err != nil {
		ErrorServer(err, c)
		return
	}

	if domainModel == nil || !domainModel.IsActive {
		ErrorNotFound(errors.New("domain is not approved yet"), c)
		return
	}

	err, urlContent := models.GetURLContentByURL(cleanedUrl, dbi, false)

	if err != nil {
		ErrorServer(err, c)
		return
	}

	if urlContent == nil {
		// url is not registered yet

		urlContent = &models.URLContent{}
	}

	Success(urlContent, c)
}
