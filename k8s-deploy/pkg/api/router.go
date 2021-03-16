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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// swag API
	_ "github.com/FederatedAI/KubeFATE/k8s-deploy/docs"
)

// APIVersion API version
const APIVersion string = "v1"

// Router of API
func Router(r *gin.Engine) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "kubefate run Success",
		})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(400, gin.H{"error": "error path"})
		})

		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		cluster := new(Cluster)
		cluster.Router(v1)

		user := new(User)
		user.Router(v1)

		job := new(Job)
		job.Router(v1)

		chart := new(Chart)
		chart.Router(v1)

		version := new(Version)
		version.Router(v1)

		namespace := new(Namespace)
		namespace.Router(v1)

		kubeLog := new(kubeLog)
		kubeLog.Router(v1)
	}
}
