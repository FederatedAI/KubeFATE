// job
package job

import (
	"fmt"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/db"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/rs/zerolog/log"
	"time"
)

type ClusterArgs struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	ChartName    string `json:"chart_name"`
	ChartVersion string `json:"chart_version"`
	Cover        bool   `json:"cover"`
	Data         []byte `json:"data"`
}

func ClusterInstall(clusterArgs *ClusterArgs, creator string) (*db.Job, error) {
	cluster := new(db.Cluster)
	if ok := cluster.IsExisted(clusterArgs.Name, clusterArgs.Namespace); ok {
		return nil, fmt.Errorf("name=%s cluster is exited", clusterArgs.Name)
	}

	job := db.NewJob("ClusterInstall", creator)

	err := setJobByUuid(job)
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("setJobByUuid error")
		return nil, err
	}
	//  save job to db
	_, err = db.Save(job)
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("create a new job of ClusterInstall")

	//do job
	go func() {
		job.Status = db.Running_j

		if clusterArgs.Cover {
			_ = helmClean(clusterArgs.Name, clusterArgs.Namespace)
		}

		//create a cluster use parameter
		cluster := db.NewCluster(clusterArgs.Name, clusterArgs.Namespace, clusterArgs.ChartName, clusterArgs.ChartVersion)
		job.ClusterId = cluster.Uuid

		err := setJobByClusterId(job)

		if job.Status == db.Running_j {
			cluster.Status = db.Creating_c

			_, err = db.Save(cluster)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("save cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("create cluster success")
		}

		err = db.UpdateByUUID(job, job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}

		err = install(cluster, clusterArgs.Data)
		if err != nil {
			job.Result = err.Error()
			job.Status = db.Failed_j
			log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("helm install cluster error")
		} else {
			job.Result = "Cluster install success"
			job.Status = db.Running_j
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("helm install cluster success")
		}

		//todo save cluster to db
		if job.Status == db.Success_j {
			cluster.Status = db.Creating_c
			err = db.UpdateByUUID(cluster, job.ClusterId)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
			}
			log.Debug().Str("cluster status", cluster.Status.String()).Str("cluster uuid", cluster.Uuid).Msg("update cluster success")
		}

		// todo job start status stop timeout
		for job.Status == db.Running_j {

			if job.TimeOut() {
				job.Result = "Checkout cluster status timeOut!"
				job.Status = db.Failed_j
				break
			}

			clusterStatusOk, err := service.CheckClusterStatus(clusterArgs.Name, clusterArgs.Namespace)
			if err != nil {
				job.Result = "CheckClusterStatus error:" + err.Error()
				job.Status = db.Failed_j
				break
			}
			if clusterStatusOk {
				job.Status = db.Success_j
				break
			}
			time.Sleep(5 * time.Second)
		}

		if job.Status == db.Canceled_j {
			job.Result = "Job canceled"
		}

		//todo save cluster to db
		if job.Status == db.Success_j {
			cluster.Status = db.Running_c
			cluster.Revision += 1
			err = db.UpdateByUUID(cluster, job.ClusterId)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
			}
			log.Debug().Str("chart_version", cluster.ChartVersion).Str("status", cluster.Status.String()).Str("cluster uuid", cluster.Uuid).Msg("update cluster success")
		}

		// rollBACK
		if job.Status != db.Success_j && job.Status != db.Canceled_j {
			err = db.ClusterDeleteByUUID(job.ClusterId)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("delete cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("delete cluster success")
		}

		job.EndTime = time.Now()
		err = db.UpdateByUUID(job, job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}

		deleteJob(job)

		if job.Status == db.Success_j {
			log.Debug().Interface("job", job).Msg("job run success")
		} else {
			log.Warn().Interface("job", job).Msg("job run failed")
		}

	}()

	return job, nil
}

func ClusterUpdate(clusterArgs *ClusterArgs, creator string) (*db.Job, error) {
	//cluster := new(db.Cluster)
	//if ok := cluster.IsExisted(clusterArgs.Name, clusterArgs.Namespace); ok {
	//	return nil, fmt.Errorf("name=%s cluster is not exited",clusterArgs.Name)
	//}
	job := db.NewJob("ClusterUpdate", creator)
	//create a cluster use parameter
	cluster, err := db.ClusterFindByName(clusterArgs.Name, clusterArgs.Namespace)
	if err != nil {
		log.Error().Err(err).Interface("clusterArgs", clusterArgs).Msg("find cluster by clusterArgs error, cluster is not exited")
		return nil, err
	}

	cluster_old := cluster

	job.ClusterId = cluster.Uuid

	err = setJobByUuid(job)
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("setJobByUuid error")
		return nil, err
	}
	err = setJobByClusterId(job)
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("setJobByUuid error")
		return nil, err
	}
	//  save job to db
	_, err = db.Save(job)
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("create a new job of ClusterUpdate")

	//do job
	go func() {
		job.Status = db.Running_j
		err = db.UpdateByUUID(job, job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}

		cluster.Status = db.Updating_c

		err = db.UpdateByUUID(cluster, cluster.Uuid)
		if err != nil {
			log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
		}
		log.Debug().Str("cluster uuid", cluster.Uuid).Msg("update cluster success")

		err := upgrade(cluster, clusterArgs)
		if err != nil {
			job.Result = err.Error()
			job.Status = db.Failed_j
			log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("helm upgrade cluster error")
		} else {
			job.Result = "cluster update success"
			job.Status = db.Running_j
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("helm upgrade cluster success")
		}

		// todo job start status stop timeout
		for job.Status == db.Running_j {
			if job.TimeOut() {
				job.Result = "checkout cluster status timeOut!"
				job.Status = db.Timeout_j
				break
			}

			clusterStatusOk, err := service.CheckClusterStatus(clusterArgs.Name, clusterArgs.Namespace)
			if err != nil {
				job.Result = "CheckClusterStatus error:" + err.Error()
				job.Status = db.Failed_j
				break
			}
			if clusterStatusOk {
				job.Status = db.Success_j
				break
			}
			time.Sleep(5 * time.Second)
		}

		if job.Status == db.Canceled_j {
			job.Result = "Job canceled"
		}

		// save cluster to db
		if job.Status == db.Success_j {
			cluster.Status = db.Running_c
			//cluster.
			cluster.Revision += 1
			err = db.UpdateByUUID(cluster, job.ClusterId)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("update cluster success")
		}

		// rollBACK
		if job.Status != db.Success_j && job.Status != db.Canceled_j {
			//todo helm rollBack

			err = db.UpdateByUUID(cluster_old, cluster_old.Uuid)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("rollBACK cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("rollBACK cluster success")
		}

		job.EndTime = time.Now()
		err = db.UpdateByUUID(job, job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}
		deleteJob(job)
		if job.Status == db.Success_j {
			log.Debug().Interface("job", job).Msg("job run success")
		} else {
			log.Warn().Interface("job", job).Msg("job run failed")
		}
	}()

	return job, nil
}

func ClusterDelete(clusterId string, creator string) (*db.Job, error) {

	job := db.NewJob("ClusterDelete", creator)
	cluster, err := db.ClusterFindByUUID(clusterId)
	if err != nil {
		log.Error().Err(err).Interface("clusterId", clusterId).Msg("find cluster by clusterId error")
		return nil, err
	}

	//cluster_old := cluster

	job.ClusterId = cluster.Uuid

	err = setJobByUuid(job)
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("setJobByUuid error")
		return nil, err
	}

	if ok := IsExistedJobByClusterID(job); ok {
		jobOther := getJobByClusterID(cluster.Uuid)
		//Cancel other job
		jobOther.Status = db.Canceled_j
	}

	err = setJobByClusterId(job)
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("setJobByClusterId error")
		return nil, err
	}

	// save job to db
	_, err = db.Save(job)
	if err != nil {
		log.Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("create a new job of ClusterDelete")

	go func() {
		job.Status = db.Running_j
		err = db.UpdateByUUID(job, job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}
		cluster.Status = db.Deleting_c
		err = db.UpdateByUUID(cluster, job.ClusterId)
		if err != nil {
			log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
		}
		log.Debug().Str("cluster uuid", cluster.Uuid).Msg("update cluster success")

		err = uninstall(cluster)
		if err != nil {
			job.Result = err.Error()
			job.Status = db.Failed_j
			log.Err(err).Str("ClusterId", cluster.Uuid).Msg("helm delete cluster error")
		} else {
			job.Result = "uninstall success"
			job.Status = db.Running_j
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("helm delete cluster success")
		}

		if job.Status == db.Running_j {
			job.Status = db.Success_j
		}

		if job.Status == db.Canceled_j {
			job.Result = "Job canceled"
		}

		//if job.Status == db.Success_j {
		err = db.ClusterDeleteByUUID(clusterId)
		if err != nil {
			log.Err(err).Interface("cluster", cluster).Msg("db delete cluster error")
		}
		log.Debug().Str("clusterUuid", clusterId).Msg("db delete cluster success")
		//}

		//// rollBACK
		//if job.Status == db.Failed_j {
		//	cluster.Status = db.Running_c
		//	err = db.UpdateByUUID(cluster_old, job.ClusterId)
		//	if err != nil {
		//		log.Error().Err(err).Interface("cluster", cluster).Msg("rollBACK cluster error")
		//	}
		//	log.Debug().Str("cluster uuid", cluster.Uuid).Msg("rollBACK cluster success")
		//}

		job.EndTime = time.Now()
		err = db.UpdateByUUID(job, job.Uuid)
		if err != nil {
			log.Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}
		deleteJob(job)
		if job.Status == db.Success_j {
			log.Debug().Interface("job", job).Msg("job run success")
		} else {
			log.Warn().Interface("job", job).Msg("job run failed")
		}
	}()

	return job, nil
}

func install(fc *db.Cluster, values []byte) error {

	err := service.RepoAddAndUpdate()
	if err != nil {
		log.Warn().Err(err).Msg("RepoAddAndUpdate error, check kubefate.yaml at env FATECLOUD_REPO_URL values,")
	}
	v := new(service.Value)
	v.Val = values
	v.T = "json"

	fc.Values = string(values)

	if fc.ChartName == "" {
		fc.ChartName = "fate"
	}

	result, err := service.Install(fc.NameSpace, fc.Name, fc.ChartName, fc.ChartVersion, v)
	if err != nil {
		return err
	}
	log.Debug().Interface("result", result).Msg("service.Install got ")

	fc.ChartName = result.ChartName
	fc.NameSpace = result.Namespace
	fc.ChartVersion = result.ChartVersion
	fc.ChartValues = result.ChartValues
	fc.Config = result.Config

	return nil
}

func upgrade(fc *db.Cluster, clusterArgs *ClusterArgs) error {

	err := uninstall(fc)
	if err != nil {
		return err
	}

	fc.ChartName = clusterArgs.ChartName
	fc.ChartVersion = clusterArgs.ChartVersion

	err = install(fc, clusterArgs.Data)
	if err != nil {
		return err
	}
	return nil
}
func uninstall(fc *db.Cluster) error {

	_, err := service.Delete(fc.NameSpace, fc.Name)

	return err
}

func helmClean(NameSpace, Name string) error {
	_, err := service.Delete(NameSpace, Name)
	return err
}

type Job interface {
	save() error
	doJob() error
	checkStatus() error
	update() error
}

func Run(j Job) (*db.Job, error) {
	var clusterArgs *ClusterArgs
	//if ok := service.IsExited(clusterArgs.Name, clusterArgs.Namespace); !ok {
	//	return nil, errors.New("cluster is exited")
	//}
	job := db.NewJob("ClusterInstall", "")

	//  save job to db
	_, err := db.Save(job)
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("create a new job of ClusterInstall")

	go func() {
		//create a cluster use parameter
		cluster := db.NewCluster(clusterArgs.Name, clusterArgs.Namespace, clusterArgs.ChartName, clusterArgs.ChartVersion)
		job.ClusterId = cluster.Uuid

		err := install(cluster, clusterArgs.Data)
		if err != nil {
			job.Result = err.Error()
			job.Status = db.Failed_j
			log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("install cluster error")
		} else {
			job.Result = "install success"
			job.Status = db.Running_j
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("install cluster success")
		}

		// todo job start status stop timeout
		for job.Status == db.Running_j {
			clusterStatusOk, err := service.CheckClusterStatus(clusterArgs.Name, clusterArgs.Namespace)
			if err != nil {
				job.Status = db.Failed_j
				break
			}
			if clusterStatusOk {
				job.Status = db.Success_j
				break
			}
			time.Sleep(time.Second)
		}

		//todo save cluster to db
		if job.Status == db.Success_j {
			_, err = db.Save(cluster)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("save cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("create cluster success")
		}

		job.EndTime = time.Now()
		err = db.UpdateByUUID(job, job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}

		log.Debug().Interface("job", job).Msg("job run success")
	}()

	return job, nil
}
