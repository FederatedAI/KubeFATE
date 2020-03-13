package job

import (
	"errors"
	"fate-cloud-agent/pkg/db"
	"sync"
)

type GlobalJobs struct {
	JobByUuid      map[string]*db.Job
	JobByClusterId map[string]*db.Job

	CS sync.Mutex
}

var GlobalJobList = initGlobalJob()

//func init() {
//	log.Debug().Interface("GlobalJobList", GlobalJobList).Msg("init")
//}

func initGlobalJob() *GlobalJobs {
	//log.Debug().Msg("initGlobalJob success")
	return &GlobalJobs{
		JobByUuid:      make(map[string]*db.Job),
		JobByClusterId: make(map[string]*db.Job),
		CS:             sync.Mutex{},
	}
}

func getJobByUUID(uuid string) *db.Job {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	return GlobalJobList.JobByUuid[uuid]
}

func getJobByClusterID(id string) *db.Job {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	return GlobalJobList.JobByClusterId[id]
}

func setJobByClusterId(job *db.Job) error {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	if _, ok := GlobalJobList.JobByUuid[job.ClusterId]; ok {
		return errors.New("cluster job is existed")
	}
	GlobalJobList.JobByClusterId[job.ClusterId] = job
	return nil
}
func setJobByUuid(job *db.Job) error {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()

	if _, ok := GlobalJobList.JobByUuid[job.Uuid]; ok {
		return errors.New("uuid job  is existed")
	}
	GlobalJobList.JobByUuid[job.Uuid] = job
	return nil
}

func IsExistedJobByClusterID(job *db.Job) bool {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	_, ok := GlobalJobList.JobByUuid[job.ClusterId]
	return ok
}

func IsExistedJobByUuid(job *db.Job) bool {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	_, ok := GlobalJobList.JobByUuid[job.Uuid]
	return ok
}

func deleteJob(job *db.Job) {
	GlobalJobList.CS.Lock()
	defer GlobalJobList.CS.Unlock()
	delete(GlobalJobList.JobByClusterId, job.ClusterId)
	delete(GlobalJobList.JobByUuid, job.Uuid)
}
