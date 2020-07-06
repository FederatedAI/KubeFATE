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
	"errors"
)

func (e *User) DropTable() {
	db.DropTable(&User{})
}

func (e *User) InitTable() {
	db.AutoMigrate(&User{})
}

func (e *User) GetList() ([]User, error) {

	var users Users
	table := db.Model(e)
	if e.Uuid != "" {
		table = table.Where("user_id = ?", e.Uuid)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.Salt != "" {
		table = table.Where("salt = ?", e.Salt)
	}

	if e.Email != "" {
		table = table.Where("email = ?", e.Email)
	}

	if e.Status != 0 {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (e *User) Get() (User, error) {

	var user User
	table := db.Model(e)
	if e.Uuid != "" {
		table = table.Where("user_id = ?", e.Uuid)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.Salt != "" {
		table = table.Where("salt = ?", e.Salt)
	}

	if e.Email != "" {
		table = table.Where("email = ?", e.Email)
	}

	if e.Status != 0 {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (e *User) Insert() (id int, err error) {
	if err = e.Encrypt(); err != nil {
		return
	}

	// check username
	var count int
	db.Model(&User{}).Where("username = ?", e.Username).Count(&count)
	if count > 0 {
		err = errors.New("Account already exists!")
		return
	}

	//Add data
	if err = db.Model(&User{}).Create(&e).Error; err != nil {
		return
	}
	id = int(e.ID)
	return
}

func (e *User) Update(id int) (update User, err error) {
	if err = db.First(&update, id).Error; err != nil {
		return
	}

	if err = db.Model(&update).Updates(&e).Error; err != nil {
		return
	}

	if err != nil {
		return
	}
	return
}

func (e *User) Delete(id int) (success bool, err error) {
	if err = db.Where("ID = ?", id).Delete(e).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}
