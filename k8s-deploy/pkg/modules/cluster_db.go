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

func (e *Cluster) DropTable() {
	db.DropTable(&Cluster{})
}

func (e *Cluster) InitTable() {
	db.AutoMigrate(&Cluster{})
}

func (e *Cluster) GetList() ([]Cluster, error) {

	var clusters Clusters
	table := db.Model(e)
	if e.Uuid != "" {
		table = table.Where("uuid = ?", e.Uuid)
	}

	if e.Name != "" {
		table = table.Where("name = ?", e.Name)
	}

	if e.NameSpace != "" {
		table = table.Where("name_space = ?", e.NameSpace)
	}

	if e.ChartName != "" {
		table = table.Where("chart_name = ?", e.ChartName)
	}

	if e.ChartVersion != "" {
		table = table.Where("chart_version = ?", e.ChartVersion)
	}

	if e.Status != 0 {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.Find(&clusters).Error; err != nil {
		return nil, err
	}
	return clusters, nil
}

func (e *Cluster) Get() (Cluster, error) {

	var cluster Cluster
	table := db.Model(e)
	if e.Uuid != "" {
		table = table.Where("uuid = ?", e.Uuid)
	}

	if e.Name != "" {
		table = table.Where("name = ?", e.Name)
	}

	if e.NameSpace != "" {
		table = table.Where("name_space = ?", e.NameSpace)
	}

	if e.ChartName != "" {
		table = table.Where("chart_name = ?", e.ChartName)
	}

	if e.ChartVersion != "" {
		table = table.Where("chart_version = ?", e.ChartVersion)
	}

	if e.Status != 0 {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.First(&cluster).Error; err != nil {
		return Cluster{}, err
	}
	return cluster, nil
}

func (e *Cluster) Insert() (id int, err error) {

	// check name namespace
	var count int
	db.Model(&Cluster{}).Where("name = ?", e.Name).Where("name_space = ?", e.NameSpace).Count(&count)
	if count > 0 {
		err = errors.New("account already exists")
		return
	}

	//Add data
	if err = db.Model(&Cluster{}).Create(&e).Error; err != nil {
		return
	}
	id = int(e.ID)
	return
}

func (e *Cluster) Update(id int) (update Cluster, err error) {
	if err = db.First(&update, id).Error; err != nil {
		return
	}

	if err = db.Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *Cluster) Delete(id int) (success bool, err error) {
	if err = db.Where("ID = ?", id).Delete(e).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}
