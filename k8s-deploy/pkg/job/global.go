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
package job

import (
	"errors"
	"sync"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
)

type GlobalJobs struct {
	JobByUuid      map[string]*modules.Job
	JobByClusterId map[string]*modules.Job

	CS sync.Mutex
}

var GlobalJobList = initGlobalJob()

//func init() {
//	log.Debug().Interface("GlobalJobList", GlobalJobList).Msg("init")
//}

func initGlobalJob() *GlobalJobs {
	//log.Debug().Msg("initGlobalJob success")
	return &GlobalJobs{
		JobByUuid:      make(map[string]*modules.Job),
		JobByClusterId: make(map[string]*modules.Job),
		CS:             sync.Mutex{},
	}
}

func getJobByUUID(uuid string) *modules.Job {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	return GlobalJobList.JobByUuid[uuid]
}

func getJobByClusterID(id string) *modules.Job {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	return GlobalJobList.JobByClusterId[id]
}

func setJobByClusterId(job *modules.Job) error {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	if _, ok := GlobalJobList.JobByUuid[job.ClusterId]; ok {
		return errors.New("cluster job is existed")
	}
	GlobalJobList.JobByClusterId[job.ClusterId] = job
	return nil
}
func setJobByUuid(job *modules.Job) error {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()

	if _, ok := GlobalJobList.JobByUuid[job.Uuid]; ok {
		return errors.New("uuid job  is existed")
	}
	GlobalJobList.JobByUuid[job.Uuid] = job
	return nil
}

func IsExistedJobByClusterID(job *modules.Job) bool {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	_, ok := GlobalJobList.JobByUuid[job.ClusterId]
	return ok
}

func IsExistedJobByUuid(job *modules.Job) bool {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	_, ok := GlobalJobList.JobByUuid[job.Uuid]
	return ok
}

func deleteJob(job *modules.Job) {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	delete(GlobalJobList.JobByClusterId, job.ClusterId)
	delete(GlobalJobList.JobByUuid, job.Uuid)
}
