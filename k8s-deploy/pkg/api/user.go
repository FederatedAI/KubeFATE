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
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// User model
type User modules.User

// Router is user router definition method
func (u *User) Router(r *gin.RouterGroup) {

	authMiddleware, _ := GetAuthMiddleware()
	user := r.Group("/user")
	{
		user.POST("/login", authMiddleware.LoginHandler)
		user.POST("/logout", authMiddleware.LogoutHandler)

		//user.GET("/findByName",u.findUser)
		//user.GET("/findByStatus",u.findUser)
	}
	user.Use(authMiddleware.MiddlewareFunc())
	{
		user.POST("", u.createUser)
		user.GET("/:userId", u.getUser)
		user.PUT("/:userId", u.setUser)
		user.DELETE("/:userId", u.deleteUser)
	}
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
		log.Info().Str("Username", u.Username).Bool("ok", ok).Msg("user delete success")
	}

	u = modules.NewUser(username, password, "")
	_, err := u.Insert()
	if err != nil {
		log.Err(err).Str("userName", username).Msg("user save error")
		return err
	}
	log.Info().Str("userUuid", u.Uuid).Str("userName", username).Msg("user  save success")
	return nil
}

func (*User) createUser(c *gin.Context) {

	user := new(modules.User)
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Use db.Newuser method to generate uuid
	user = modules.NewUser(user.Username, user.Password, user.Email)
	if user.IsExisted() {
		c.JSON(500, gin.H{"error": USEREXISTED})
		return
	}
	_, err := user.Insert()
	if err != nil {
		c.JSON(500, gin.H{"msg": err})
	}

	c.JSON(200, gin.H{"msg": "createCluster success", "data": user})
}

func (*User) setUser(c *gin.Context) {

	userId := c.Param("userId")
	user := new(modules.User)
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user.Uuid = userId
	_, err := user.Update(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	c.JSON(200, gin.H{"msg": "setUser success"})
}
func (*User) getUser(c *gin.Context) {

	userId := c.Param("userId")
	if userId == "" {
		c.JSON(400, gin.H{"error": "err"})
	}
	result, err := getUserFindByUUID(userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	c.JSON(200, gin.H{"data": result})
}

func getUserFindByUUID(uuid string) (modules.User, error) {
	u := modules.User{Uuid: uuid}
	user, err := u.Get()
	return user, err
}

func (*User) deleteUser(c *gin.Context) {

	userId := c.Param("userId")
	if userId == "" {
		c.JSON(400, gin.H{"error": "err"})
	}
	u := modules.User{Uuid: userId}
	_, err := u.Delete()
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	c.JSON(200, gin.H{"msg": "deleteUser success"})
}
