/*
 *  Copyright 2019-2020 VMware, Inc.
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
 */

package modules

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/orm"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/pbkdf2"
	"k8s.io/apimachinery/pkg/util/rand"
)

const saltSize = 128

type User struct {
	Uuid     string     `json:"uuid,omitempty"  gorm:"index:uuid;unique"`
	Username string     `json:"username,omitempty" gorm:"index:username;unique;not null"`
	Password string     `json:"password,omitempty" gorm:"type:varchar(512)"`
	Salt     string     `json:"salt,omitempty"`
	Email    string     `json:"email,omitempty"`
	Status   UserStatus `json:"userStatus,omitempty"`
	gorm.Model
}

type Users []User

type UserStatus int8

const (
	Deprecate_u UserStatus = iota
	Available_u
)

func (s UserStatus) String() string {
	names := []string{
		"Deprecate",
		"Available",
	}

	return names[s]
}

func (s UserStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func encryption(plaintext, salt string) string {
	iterations := 100000
	digest := sha256.New
	secretaries := pbkdf2.Key([]byte(plaintext), []byte(salt), iterations, 256, digest)
	return fmt.Sprintf("%x", secretaries)
}

func NewUser(username string, password string, email string) *User {
	salt := rand.String(saltSize)
	u := &User{
		Uuid:     uuid.NewV4().String(),
		Username: username,
		Password: encryption(password, salt),
		Salt:     salt,
		Email:    email,
		Status:   Deprecate_u,
	}

	return u
}

func (e *User) IsValid() bool {
	db := orm.DBCLIENT
	gotUser := new(User)
	db.Where("username = ?", e.Username).First(&gotUser)
	if db.Error != nil || gotUser == nil {
		return false
	}

	if gotUser.Password != encryption(e.Password, gotUser.Salt) {
		return false
	}
	return true
}

func (e *User) IsExisted() bool {

	gotUsers := make(Users, 0)
	db.Where("username = ?", e.Username).First(&gotUsers)
	if db.Error != nil || len(gotUsers) == 0 {
		return false
	}
	return true
}

func (e *User) Encrypt()  (err error)  {
	salt := rand.String(saltSize)
	e.Password = encryption(e.Password,salt)
	return
}

