/*
 * Copyright 2019-2020 VMware, Inc.
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
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"k8s.io/apimachinery/pkg/util/rand"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var identityKey = "id"

// AuthMiddleware singleton middleware
var AuthMiddleware *jwt.GinJWTMiddleware = nil

// GetAuthMiddleware get auth middleware in other file.
func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	var err error = nil
	if AuthMiddleware == nil {
		err = initAuthmiddleware()
	}

	return AuthMiddleware, err
}

func generateKey() []byte {
	return []byte(rand.String(32))
}

func initAuthmiddleware() error {
	tmpAuth, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "gin jwt",
		Key:         generateKey(),
		Timeout:     time.Minute * 30,
		MaxRefresh:  time.Minute * 30,
		IdentityKey: identityKey,
		// Use username as identity key and set it as jwt token
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*modules.User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},

		// Get identity key from jwt token, convert to user
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Username: claims[identityKey].(string),
			}
		},

		// Check if user exists
		Authenticator: func(c *gin.Context) (interface{}, error) {
			loginVals := new(modules.User)
			if err := c.ShouldBindJSON(loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			log.Debug().Msg("Login user name: " + loginVals.Username)
			if loginVals.IsValid() {
				return loginVals, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},

		// Check if the user has relevant permissions
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// if v, ok := data.(*User); ok && v.Username == "admin" {
			if v, ok := data.(*User); ok && v.Username != "" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	AuthMiddleware = tmpAuth
	return err
}
