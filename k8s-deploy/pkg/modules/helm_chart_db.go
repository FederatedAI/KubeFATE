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

func (e *HelmChart) DropTable() {
	DB.DropTable(&HelmChart{})
}

func (e *HelmChart) InitTable() {
	DB.AutoMigrate(&HelmChart{})
}

func (e *HelmChart) GetList() ([]HelmChart, error) {

	var helmCharts HelmCharts
	table := DB.Model(e)
	if e.Uuid != "" {
		table = table.Where("uuid = ?", e.Uuid)
	}

	if e.Name != "" {
		table = table.Where("name = ?", e.Name)
	}

	if e.Chart != "" {
		table = table.Where("chart = ?", e.Chart)
	}

	if e.Version != "" {
		table = table.Where("version = ?", e.Version)
	}

	if e.AppVersion != "" {
		table = table.Where("app_version = ?", e.AppVersion)
	}

	if err := table.Find(&helmCharts).Error; err != nil {
		return nil, err
	}
	return helmCharts, nil
}

func (e *HelmChart) Get() (HelmChart, error) {

	var cluster HelmChart
	table := DB.Model(e)
	if e.Uuid != "" {
		table = table.Where("uuid = ?", e.Uuid)
	}

	if e.Name != "" {
		table = table.Where("name = ?", e.Name)
	}

	if e.Chart != "" {
		table = table.Where("chart = ?", e.Chart)
	}

	if e.Version != "" {
		table = table.Where("version = ?", e.Version)
	}

	if e.AppVersion != "" {
		table = table.Where("app_version = ?", e.AppVersion)
	}

	if err := table.First(&cluster).Error; err != nil {
		return HelmChart{}, err
	}
	return cluster, nil
}

func (e *HelmChart) Insert() (id int, err error) {

	// check name namespace
	var count int
	DB.Model(&HelmChart{}).Where("version = ?", e.Version).Count(&count)
	if count > 0 {
		err = errors.New("helmChart already exists, version = " + e.Version)
		return
	}

	//Add data
	if err = DB.Model(&HelmChart{}).Create(&e).Error; err != nil {
		return
	}
	id = int(e.ID)
	return
}

func (e *HelmChart) Upload() (err error) {
	var count int
	if err = DB.Model(&HelmChart{}).Where("version = ?", e.Version).Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		DB.Delete(&HelmChart{}, "version = ?", e.Version)
	}
	if err = DB.Model(&HelmChart{}).Create(&e).Error; err != nil {
		return
	}
	return
}

func (e *HelmChart) Update(id int) (update HelmChart, err error) {
	if err = DB.First(&update, id).Error; err != nil {
		return
	}

	if err = DB.Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *HelmChart) Delete(id int) (success bool, err error) {
	if err = DB.Where("ID = ?", id).Delete(e).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (e *HelmChart) DeleteByUuid(Uuid string) (success bool, err error) {
	if err = DB.Where("uuid = ?", Uuid).Delete(e).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (e *HelmChart) IsExisted() bool {
	var count int
	DB.Model(&HelmChart{}).Where("name = ?", e.Name).Where("version = ?", e.Version).Count(&count)
	if DB.Error == nil && count > 0 {
		return true
	}
	return false
}
