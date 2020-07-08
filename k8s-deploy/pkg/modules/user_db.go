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
	"errors"
)

func (e *User) DropTable() {
	DB.DropTable(&User{})
}

func (e *User) InitTable() {
	DB.AutoMigrate(&User{})
}

func (e *User) GetList() ([]User, error) {

	var users Users
	table := DB.Model(e)
	if e.Uuid != "" {
		table = table.Where("uuid = ?", e.Uuid)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
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

	for _, u := range users {
		u.Password = ""
	}

	return users, nil
}

func (e *User) Get() (User, error) {

	var user User
	table := DB.Model(e).Unscoped()
	if e.Uuid != "" {
		table = table.Where("uuid = ?", e.Uuid)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
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

	user.Password = ""

	return user, nil
}

func (e *User) Insert() (id int, err error) {
	e.Encrypt()

	// check username
	var count int
	DB.Model(&User{}).Where("username = ?", e.Username).Count(&count)
	if count > 0 {
		err = errors.New("Account already exists!")
		return
	}

	//Add data
	if err = DB.Model(&User{}).Create(&e).Error; err != nil {
		return
	}
	id = int(e.ID)
	return
}

func (e *User) Update(id uint) (update User, err error) {
	if err = DB.First(&update, id).Error; err != nil {
		return
	}

	e.Encrypt()

	if err = DB.Model(&update).Updates(&e).Error; err != nil {
		return
	}

	update.Password = ""

	return
}

func (e *User) DeleteById(id uint) (success bool, err error) {
	if err = DB.Unscoped().Where("ID = ?", id).Delete(e).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (e *User) Delete() (success bool, err error) {
	user, err := e.Get()
	if err != nil {
		return false, err
	}
	return user.DeleteById(user.ID)
}
