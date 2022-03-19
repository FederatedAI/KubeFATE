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
	"errors"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// User model
type User modules.User

// Router is user router definition method
func (u *User) Router(r *gin.RouterGroup) {
	authM, _ := GetAuthMiddleware()
	auth := authMiddleware{authM}

	user := r.Group("/user")
	{
		user.POST("/login", auth.login)
		user.POST("/logout", auth.logout)

		//user.GET("/findByName",u.findUser)
		//user.GET("/findByStatus",u.findUser)
	}
	user.Use(auth.MiddlewareFunc())
	{
		user.POST("", u.createUser)
		user.GET("/", u.getUserList)
		user.GET("/:userId", u.getUser)
		user.PUT("/:userId", u.setUser)
		user.DELETE("/:userId", u.deleteUser)
	}
}

type authMiddleware struct {
	*jwt.GinJWTMiddleware
}

// login login
// @Summary login
// @Tags User
// @Accept  json
// @Produce  json
// @Param login body Login true "Login"
// @Success 200 {object} TokenResult
// @Failure 401 {object} JSONEMSGResult
// @Router /user/login [post]
func (authMiddleware *authMiddleware) login(c *gin.Context) {
	authMiddleware.LoginHandler(c)
}

// logout logout
// @Summary logout
// @Tags User
// @Produce  json
// @Success 200 {object} string
// @Failure 401 {object} JSONEMSGResult
// @Router /user/logout [post]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (authMiddleware *authMiddleware) logout(c *gin.Context) {
	authMiddleware.LogoutHandler(c)
}

func generateAdminUser() error {
	username := viper.GetString("user.username")
	password := viper.GetString("user.password")

	u := modules.NewUser(username, password, "")
	if u.IsExisted() {
		user := modules.User{Username: username}
		ok, err := user.Delete()
		if err != nil {
			log.Err(err).Str("userName", username).Msg("Delete user by name error")
			return err
		}
		log.Info().Str("Username", u.Username).Bool("ok", ok).Msg("user delete Success")
	}

	u = modules.NewUser(username, password, "")
	_, err := u.Insert()
	if err != nil {
		log.Err(err).Str("userName", username).Msg("user save error")
		return err
	}
	log.Info().Str("userUuid", u.Uuid).Str("userName", username).Msg("user  save Success")
	return nil
}

// createUser Create a user
// @Summary Create a user
// @Tags User
// @Produce  json
// @Param  User body modules.User true "User"
// @Success 200 {object} JSONResult{data=modules.User} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /user [post]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*User) createUser(c *gin.Context) {

	user := new(modules.User)
	if err := c.ShouldBindJSON(user); err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Use db.Newuser method to generate uuid
	user = modules.NewUser(user.Username, user.Password, user.Email)
	if user.IsExisted() {
		log.Error().Err(errors.New(USEREXISTED)).Msg("request error")
		c.JSON(500, gin.H{"error": USEREXISTED})
		return
	}
	_, err := user.Insert()
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"msg": err.Error()})
	}

	log.Debug().Interface("data", "user").Msg("result")
	c.JSON(200, gin.H{"msg": "createCluster Success", "data": user})
}

// setUser Update user
// @Summary Update user
// @Tags User
// @Produce  json
// @Param  User body modules.User true "User"
// @Success 200 {object} JSONResult{data=modules.User} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /user [put]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*User) setUser(c *gin.Context) {

	userID := c.Param("userId")
	user := new(modules.User)
	if err := c.ShouldBindJSON(user); err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user.Uuid = userID
	_, err := user.Update(user.ID)
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
	}
	log.Debug().Interface("result", "setUser Success").Msg("result")
	c.JSON(200, gin.H{"msg": "setUser Success"})
}

// getUser Get user by userId
// @Summary Get user by userId
// @Tags User
// @Produce  json
// @Param  userId path string true "User"
// @Success 200 {object} JSONResult{data=modules.User} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /user/{userId} [get]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*User) getUser(c *gin.Context) {

	userID := c.Param("userId")
	if userID == "" {
		log.Error().Err(errors.New("not exit userId")).Msg("request error")
		c.JSON(400, gin.H{"error": "not exit userId"})
	}
	result, err := getUserFindByUUID(userID)
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
	}
	log.Debug().Interface("result", result).Msg("result")
	c.JSON(200, gin.H{"msg": "getUser Success", "data": result})
}

// getUserList List all available user
// @Summary List all available user
// @Tags User
// @Produce  json
// @Success 200 {object} JSONResult{data=[]modules.User} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /user [get]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*User) getUserList(c *gin.Context) {

	u := new(modules.User)
	result, err := u.GetList()
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
	}
	log.Debug().Interface("result", result).Msg("result")
	c.JSON(200, gin.H{"msg": "getUserList Success", "data": result})
}

func getUserFindByUUID(uuid string) (modules.User, error) {
	u := modules.User{Uuid: uuid}
	user, err := u.Get()
	return user, err
}

// deleteUser Delete user by userId
// @Summary Delete user by userId
// @Tags User
// @Produce  json
// @Param  userId path string true "User"
// @Success 200 {object} JSONEMSGResult "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /user/{userId} [delete]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*User) deleteUser(c *gin.Context) {

	userID := c.Param("userId")
	if userID == "" {
		log.Error().Err(errors.New("not exit userId")).Msg("request error")
		c.JSON(400, gin.H{"error": "not exit userId"})
	}
	u := modules.User{Uuid: userID}
	_, err := u.Delete()
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"msg": "deleteUser Success"})
}
