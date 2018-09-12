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
	"github.com/primasio/wormhole/http/middlewares"
	"github.com/primasio/wormhole/models"
	"strconv"
)

type DomainController struct{}

type DomainForm struct {
	Domain string `form:"domain" json:"domain" binding:"required"`
	Title  string `form:"title" json:"title" binding:"required"`
}

func (ctrl *DomainController) Create(c *gin.Context) {
	var form DomainForm

	if err := c.ShouldBind(&form); err != nil {
		Error(err.Error(), c)
	} else {
		dbi := db.GetDb()

		cleanedDomain := models.CleanDomain(form.Domain)

		// Check domain uniqueness
		err, check := models.GetDomainByDomainName(cleanedDomain, db.GetDb(), false)

		if err != nil {
			ErrorServer(err, c)
			return
		}

		if check != nil {
			Error("domain exists", c)
			return
		}

		// Save domain to db
		domainModel := &models.Domain{}

		domainModel.Domain = cleanedDomain
		domainModel.Title = form.Title
		domainModel.HashKey = models.GetDomainHashKey(cleanedDomain)

		userId, _ := c.Get(middlewares.AuthorizedUserId)
		domainModel.UserID = userId.(uint)

		dbi.Create(&domainModel)

		Success(domainModel, c)
	}
}

func (ctrl *DomainController) Get(c *gin.Context) {

	domain := c.Query("domain")

	if domain == "" {
		Error("missing query param domain", c)
		return
	}

	cleanedDomain := models.CleanDomain(domain)

	err, domainModel := models.GetDomainByDomainName(cleanedDomain, db.GetDb(), false)

	if err != nil {
		ErrorServer(err, c)
		return
	}

	if domainModel == nil {
		ErrorNotFound(errors.New("domain not found"), c)
		return
	}

	Success(domainModel, c)
}

func (ctrl *DomainController) List(c *gin.Context) {

	urlType := c.Query("type")

	pageSize := 20
	page, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		page = 0
	}

	offsetNum := page * pageSize

	domainList := make([]models.Domain, 0)

	dbi := db.GetDb()
	query := dbi.Where("is_active = ?", urlType != "voting")
	query = query.Order("created_at DESC").Offset(offsetNum).Limit(pageSize)

	query.Find(&domainList)

	Success(domainList, c)
}

func (ctrl *DomainController) Vote(c *gin.Context) {

	domain := c.Query("domain")

	if domain == "" {
		Error("missing query param domain", c)
		return
	}

	cleanedDomain := models.CleanDomain(domain)

	err, domainModel := models.GetDomainByDomainName(cleanedDomain, db.GetDb(), false)

	if err != nil {
		ErrorServer(err, c)
		return
	}

	if domainModel == nil {
		ErrorNotFound(errors.New("domain not found"), c)
		return
	}

	if domainModel.IsActive {
		Error("domain is already active", c)
		return
	}

	userId, _ := c.Get(middlewares.AuthorizedUserId)
	userIdNum := userId.(uint)

	if domainModel.UserID == userIdNum {
		Error("user already voted", c)
		return
	}

	// Update vote should be executed on a locked object

	lockedDomain := &models.Domain{}

	tx := db.GetDb().Begin()

	// If SQLite is used, FOR UPDATE is not supported
	// Then there is an error of concurrent votes count

	sql := "SELECT * FROM domains WHERE id = ?"

	if db.GetDbType() != db.SQLITE {
		sql = sql + " FOR UPDATE"
	}

	tx.Raw(sql, domainModel.ID).Scan(&lockedDomain)

	if lockedDomain.ID == 0 {
		tx.Rollback()
		ErrorServer(errors.New("error lock domain"), c)
		return
	}

	// Check user vote status
	// this should be performed after the locking of url_content
	// to avoid race condition of concurrent voting from the same user

	vote := &models.DomainVote{
		UserID:   userIdNum,
		DomainID: lockedDomain.ID,
	}

	tx.Where(&vote).First(&vote)

	if vote.ID != 0 {
		tx.Rollback()
		Error("user already voted", c)
		return
	}

	lockedDomain.Votes++

	if err := tx.Save(&lockedDomain).Error; err != nil {
		tx.Rollback()
		ErrorServer(err, c)
		return
	}

	if err := tx.Create(vote).Error; err != nil {
		tx.Rollback()
		ErrorServer(err, c)
		return
	}

	tx.Commit()
	Success(lockedDomain, c)
}

func (ctrl *DomainController) Approve(c *gin.Context) {
	domain := c.Query("domain")

	if domain == "" {
		Error("missing query param domain", c)
		return
	}

	err, domainModel := models.GetDomainByDomainName(domain, db.GetDb(), false)

	if err != nil {
		ErrorServer(err, c)
		return
	}

	if domainModel == nil {
		ErrorNotFound(errors.New("domain not found"), c)
		return
	}

	if domainModel.IsActive {
		Error("domain is already active", c)
		return
	}

	domainModel.IsActive = true

	dbi := db.GetDb()

	dbi.Save(domainModel)

	Success(domainModel, c)
}
