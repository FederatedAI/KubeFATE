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

func (e *Job) DropTable() {
	DB.DropTable(&Job{})
}

func (e *Job) InitTable() {
	DB.AutoMigrate(&Job{})
}

func (e *Job) GetList() ([]Job, error) {

	var jobs Jobs
	table := DB.Model(e)
	if e.Uuid != "" {
		table = table.Where("uuid = ?", e.Uuid)
	}

	if e.ClusterId != "" {
		table = table.Where("cluster_id = ?", e.ClusterId)
	}

	if e.Creator != "" {
		table = table.Where("creator = ?", e.Creator)
	}

	if e.Method != "" {
		table = table.Where("method = ?", e.Method)
	}

	if e.Status != 0 {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (e *Job) Get() (Job, error) {

	var job Job
	table := DB.Model(e)
	if e.Uuid != "" {
		table = table.Where("uuid = ?", e.Uuid)
	}

	if e.ClusterId != "" {
		table = table.Where("cluster_id = ?", e.ClusterId)
	}

	if e.Creator != "" {
		table = table.Where("creator = ?", e.Creator)
	}

	if e.Method != "" {
		table = table.Where("method = ?", e.Method)
	}

	if e.Status != 0 {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.First(&job).Error; err != nil {
		return Job{}, err
	}
	return job, nil
}

func (e *Job) Insert() (id int, err error) {

	// check name namespace
	var count int
	DB.Model(&Job{}).Where("uuid = ?", e.Uuid).Count(&count)
	if count > 0 {
		err = errors.New("job already exists, uuid = " + e.Uuid)
		return
	}

	//Add data
	if err = DB.Model(&Job{}).Create(&e).Error; err != nil {
		return
	}
	id = int(e.ID)
	return
}

func (e *Job) Update(id int) (update Job, err error) {
	if err = DB.First(&update, id).Error; err != nil {
		return
	}

	if err = DB.Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *Job) UpdateByUuid(uuid string) (update Job, err error) {
	if err = DB.Where("uuid = ?", uuid).First(&update).Error; err != nil {
		return
	}

	if err = DB.Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *Job) DeleteById(id uint) (success bool, err error) {
	if err = DB.Where("ID = ?", id).Delete(e).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (e *Job) Delete() (bool, error) {
	job, err := e.Get()
	if err != nil {
		return false, err
	}
	return e.DeleteById(job.ID)
}

func (e *Job) SetStatus(status JobStatus) error {
	if err := DB.Model(e).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (e *Job) SetResult(result string) error {
	if err := DB.Model(e).Update("result", result).Error; err != nil {
		return err
	}
	return nil
}

func (e *Job) SetSubJobs(subJobs SubJobs) error {
	if err := DB.Model(e).Update("sub_jobs", subJobs).Error; err != nil {
		return err
	}
	return nil
}

func (e *Job) IsExisted(uuid string) bool {
	var count int
	DB.Model(&Job{}).Where("uuid = ?", uuid).Count(&count)
	if DB.Error == nil && count > 0 {
		return true
	}
	return false
}
