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
	"github.com/gin-gonic/gin"
	"github.com/primasio/wormhole/db"
	"github.com/primasio/wormhole/http/middlewares"
	"github.com/primasio/wormhole/models"
)

type ArticleController struct{}

func (ctrl *ArticleController) Publish(c *gin.Context) {
	var article models.Article

	if err := c.ShouldBind(&article); err != nil {
		Error(err.Error(), c)
	} else {
		dbi := db.GetDb()

		userId, _ := c.Get(middlewares.AuthorizedUserId)

		article.UserId = userId.(uint)

		dbi.Create(&article)

		// TODO: Create async task to publish article to Primas

		Success(article, c)
	}
}

func (ctrl *ArticleController) Get(c *gin.Context) {

}
