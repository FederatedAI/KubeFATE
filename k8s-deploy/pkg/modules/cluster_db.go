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
	"fmt"
)

func (e *Cluster) DropTable() {
	DB.DropTable(&Cluster{})
}

func (e *Cluster) InitTable() {
	DB.AutoMigrate(&Cluster{})
}

func (e *Cluster) GetList() ([]Cluster, error) {

	var clusters Clusters
	table := DB.Model(e)
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

	if e.Status != 0 {
		table = table.Unscoped()
	}

	if err := table.Find(&clusters).Error; err != nil {
		return nil, err
	}
	return clusters, nil
}

func (e *Cluster) GetListAll(all bool) ([]Cluster, error) {

	var clusters Clusters
	table := DB.Model(e)
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

	if all {
		table = table.Unscoped()
	}

	if err := table.Find(&clusters).Error; err != nil {
		return nil, err
	}
	return clusters, nil
}

func (e *Cluster) Get() (Cluster, error) {

	var cluster Cluster
	table := DB.Model(e)
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
	DB.Model(&Cluster{}).Where("name = ?", e.Name).Where("name_space = ?", e.NameSpace).Count(&count)
	if count > 0 {
		err = fmt.Errorf("cluster already exists in database, name = %s, namespace=%s", e.Name, e.NameSpace)
		return
	}

	//Add data
	if err = DB.Model(&Cluster{}).Create(&e).Error; err != nil {
		return
	}
	id = int(e.ID)
	return
}

func (e *Cluster) Update(id int) (update Cluster, err error) {
	if err = DB.First(&update, id).Error; err != nil {
		return
	}

	if err = DB.Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *Cluster) UpdateByUuid(uuid string) (update Cluster, err error) {
	if err = DB.Where("uuid = ?", uuid).First(&update).Error; err != nil {
		return
	}

	if err = DB.Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *Cluster) deleteById(id uint) (success bool, err error) {
	if err = DB.Where("ID = ?", id).Delete(e).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (e *Cluster) Delete() (bool, error) {
	cluster, err := e.Get()
	if err != nil {
		return false, err
	}
	err = e.SetStatus(ClusterStatusDeleted)
	if err != nil {
		return false, err
	}
	return e.deleteById(cluster.ID)
}

func (e *Cluster) IsExisted(name, namespace string) bool {
	var count int
	DB.Model(&Cluster{}).Where("name = ?", name).Where("name_space = ?", namespace).Count(&count)
	if DB.Error == nil && count > 0 {
		return true
	}
	return false
}

func (e *Cluster) SetStatus(status ClusterStatus) error {
	if err := DB.Model(e).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (e *Cluster) SetSpec(spec MapStringInterface) error {
	if err := DB.Model(e).Update("spec", spec).Error; err != nil {
		return err
	}
	return nil
}

func (e *Cluster) SetValues(values string) error {
	if err := DB.Model(e).Update("values", values).Error; err != nil {
		return err
	}
	return nil
}
