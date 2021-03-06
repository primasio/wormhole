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
	"github.com/golang/glog"
	"net/http"
)

func Error(msg string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": msg})
}

func Success(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func ErrorServer(err error, c *gin.Context) {
	glog.Error(err)
	c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Internal Server Error"})
}

func ErrorUnauthorized(msg string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": msg})
}

func ErrorNotFound(err error, c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
}
