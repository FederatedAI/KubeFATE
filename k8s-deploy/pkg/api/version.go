/*
 * Copyright 2019-2021 VMware, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package api

import (
	"github.com/gin-gonic/gin"
)

// ServiceVersion code release version
const ServiceVersion = "v1.4.4"

// Version API struct
type Version struct {
}

// Router is Cluster router definition method
func (c *Version) Router(r *gin.RouterGroup) {

	authMiddleware, _ := GetAuthMiddleware()
	version := r.Group("/version")
	version.Use(authMiddleware.MiddlewareFunc())
	{
		version.GET("/", c.getVersion)
	}
}

// getVersion Get build version
// @Summary Get build version
// @Tags Version
// @Produce  json
// @Success 200 {object} VersionResult
// @Failure 401 {object} JSONERRORResult
// @Router /version [get]
// @Param Authorization header string true "Authentication header" default(Bearer <Token>)
// @Security ApiKeyAuth
func (*Version) getVersion(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "getVersion Success", "version": ServiceVersion})
}
