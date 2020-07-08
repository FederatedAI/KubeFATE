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
	Uuid     string     `json:"uuid,omitempty"  gorm:"type:varchar(36);index:uuid;unique"`
	Username string     `json:"username,omitempty" gorm:"index:username;unique;not null"`
	Password string     `json:"password,omitempty" gorm:"type:varchar(512);not null"`
	Salt     string     `json:"salt,omitempty" gorm:"type:varchar(128);not null"`
	Email    string     `json:"email,omitempty" gorm:"type:varchar(255)"`
	Status   UserStatus `json:"userStatus,omitempty" gorm:"size:8;not null"`
	gorm.Model
}

type Users []User

type UserStatus int8

const (
	Deprecate_u UserStatus = iota + 1
	Available_u
)

func (s UserStatus) String() string {
	names := map[UserStatus]string{
		Deprecate_u: "Deprecate",
		Available_u: "Available",
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
		Password: password,
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

	var count int
	DB.Model(&User{}).Where("username = ?", e.Username).Count(&count)
	if DB.Error == nil && count > 0 {
		return true
	}
	return false
}

func (e *User) Encrypt() {
	e.Password = encryption(e.Password, e.Salt)
	return
}
