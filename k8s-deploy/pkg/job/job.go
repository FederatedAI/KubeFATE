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

// job
package job

import (
	"fmt"
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/rs/zerolog/log"
)

type ClusterArgs struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	ChartName    string `json:"chart_name"`
	ChartVersion string `json:"chart_version"`
	Cover        bool   `json:"cover"`
	Data         []byte `json:"data"`
}

func ClusterInstall(clusterArgs *ClusterArgs, creator string) (*modules.Job, error) {
	cluster := new(modules.Cluster)
	if ok := cluster.IsExisted(clusterArgs.Name, clusterArgs.Namespace); ok {
		return nil, fmt.Errorf("name=%s cluster is exited", clusterArgs.Name)
	}

	job := modules.NewJob("ClusterInstall", creator)

	err := setJobByUuid(job)
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("setJobByUuid error")
		return nil, err
	}
	//  save job to db
	_, err = job.Insert()
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("create a new job of ClusterInstall")

	//do job
	go func() {
		job.Status = modules.Running_j

		if clusterArgs.Cover {
			log.Info().Msg("Overwrite current installation")
			err = helmClean(clusterArgs.Name, clusterArgs.Namespace)
			if err != nil {
				log.Error().Msg("helmClean error")
			}
			log.Info().Msg("HelmClean success")
		}

		//create a cluster use parameter
		cluster := modules.NewCluster(clusterArgs.Name, clusterArgs.Namespace, clusterArgs.ChartName, clusterArgs.ChartVersion)
		job.ClusterId = cluster.Uuid

		err := setJobByClusterId(job)

		if job.Status == modules.Running_j {
			cluster.Status = modules.Creating_c

			_, err = cluster.Insert()
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("save cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("create cluster success")
		}

		_, err = job.UpdateByUuid(job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}

		err = install(cluster, clusterArgs.Data)
		if err != nil {
			job.Result = err.Error()
			job.Status = modules.Failed_j
			log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("helm install cluster error")
		} else {
			job.Result = "Cluster install success"
			job.Status = modules.Running_j
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("helm install cluster success")
		}

		//todo save cluster to modules
		if job.Status == modules.Success_j {
			cluster.Status = modules.Creating_c
			_, err = cluster.UpdateByUuid(job.ClusterId)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
			}
			log.Debug().Str("cluster status", cluster.Status.String()).Str("cluster uuid", cluster.Uuid).Msg("update cluster success")
		}

		// todo job start status stop timeout
		for job.Status == modules.Running_j {

			if job.TimeOut() {
				job.Result = "Checkout cluster status timeOut!"
				job.Status = modules.Failed_j
				break
			}

			clusterStatusOk, err := service.CheckClusterStatus(clusterArgs.Name, clusterArgs.Namespace)
			if err != nil {
				job.Result = "CheckClusterStatus error:" + err.Error()
				job.Status = modules.Failed_j
				break
			}
			if clusterStatusOk {
				job.Status = modules.Success_j
				break
			}
			time.Sleep(5 * time.Second)
		}

		if job.Status == modules.Canceled_j {
			job.Result = "Job canceled"
		}

		//todo save cluster to modules
		if job.Status == modules.Success_j {
			cluster.Status = modules.Running_c
			cluster.Revision += 1
			_, err = cluster.UpdateByUuid(job.ClusterId)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
			}
			log.Debug().Str("chart_version", cluster.ChartVersion).Str("status", cluster.Status.String()).Str("cluster uuid", cluster.Uuid).Msg("update cluster success")
		}

		// rollBACK
		if job.Status != modules.Success_j && job.Status != modules.Canceled_j {
			c := &modules.Cluster{Uuid: job.ClusterId}
			_, err = c.Delete()
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("delete cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("delete cluster success")
		}

		job.EndTime = time.Now()
		_, err = job.UpdateByUuid(job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}

		deleteJob(job)

		if job.Status == modules.Success_j {
			log.Debug().Interface("job", job).Msg("job run success")
		} else {
			log.Warn().Interface("job", job).Msg("job run failed")
		}

	}()

	return job, nil
}

func ClusterUpdate(clusterArgs *ClusterArgs, creator string) (*modules.Job, error) {
	//cluster := new(modules.Cluster)
	//if ok := cluster.IsExisted(clusterArgs.Name, clusterArgs.Namespace); ok {
	//	return nil, fmt.Errorf("name=%s cluster is not exited",clusterArgs.Name)
	//}
	job := modules.NewJob("ClusterUpdate", creator)
	//create a cluster use parameter
	c := modules.Cluster{Name: clusterArgs.Name, NameSpace: clusterArgs.Namespace}
	cluster, err := c.Get()
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
	//  save job to modules
	_, err = job.Insert()
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("create a new job of ClusterUpdate")

	//do job
	go func() {
		job.Status = modules.Running_j
		_, err = job.UpdateByUuid(job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}

		cluster.Status = modules.Updating_c

		_, err = cluster.UpdateByUuid(cluster.Uuid)
		if err != nil {
			log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
		}
		log.Debug().Str("cluster uuid", cluster.Uuid).Msg("update cluster success")

		err := upgrade(&cluster, clusterArgs)
		if err != nil {
			job.Result = err.Error()
			job.Status = modules.Failed_j
			log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("helm upgrade cluster error")
		} else {
			job.Result = "cluster update success"
			job.Status = modules.Running_j
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("helm upgrade cluster success")
		}

		// todo job start status stop timeout
		for job.Status == modules.Running_j {
			if job.TimeOut() {
				job.Result = "checkout cluster status timeOut!"
				job.Status = modules.Timeout_j
				break
			}

			clusterStatusOk, err := service.CheckClusterStatus(clusterArgs.Name, clusterArgs.Namespace)
			if err != nil {
				job.Result = "CheckClusterStatus error:" + err.Error()
				job.Status = modules.Failed_j
				break
			}
			if clusterStatusOk {
				job.Status = modules.Success_j
				break
			}
			time.Sleep(5 * time.Second)
		}

		if job.Status == modules.Canceled_j {
			job.Result = "Job canceled"
		}

		// save cluster to modules
		if job.Status == modules.Success_j {
			cluster.Status = modules.Running_c
			//cluster.
			cluster.Revision += 1
			_, err = cluster.UpdateByUuid(job.ClusterId)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("update cluster success")
		}

		// rollBACK
		if job.Status != modules.Success_j && job.Status != modules.Canceled_j {
			//todo helm rollBack

			_, err = cluster_old.UpdateByUuid(cluster_old.Uuid)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("rollBACK cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("rollBACK cluster success")
		}

		job.EndTime = time.Now()
		_, err = job.UpdateByUuid(job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}
		deleteJob(job)
		if job.Status == modules.Success_j {
			log.Debug().Interface("job", job).Msg("job run success")
		} else {
			log.Warn().Interface("job", job).Msg("job run failed")
		}
	}()

	return job, nil
}

func ClusterDelete(clusterId string, creator string) (*modules.Job, error) {

	job := modules.NewJob("ClusterDelete", creator)
	c := modules.Cluster{Uuid: clusterId}
	cluster, err := c.Get()
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
		jobOther.Status = modules.Canceled_j
	}

	err = setJobByClusterId(job)
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("setJobByClusterId error")
		return nil, err
	}

	// save job to modules
	_, err = job.Insert()
	if err != nil {
		log.Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("create a new job of ClusterDelete")

	go func() {
		job.Status = modules.Running_j
		_, err = job.UpdateByUuid(job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}
		cluster.Status = modules.Deleting_c
		_, err = cluster.UpdateByUuid(job.ClusterId)
		if err != nil {
			log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
		}
		log.Debug().Str("cluster uuid", cluster.Uuid).Msg("update cluster success")

		err = uninstall(&cluster)
		if err != nil {
			job.Result = err.Error()
			job.Status = modules.Failed_j
			log.Err(err).Str("ClusterId", cluster.Uuid).Msg("helm delete cluster error")
		} else {
			job.Result = "uninstall success"
			job.Status = modules.Running_j
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("helm delete cluster success")
		}

		if job.Status == modules.Running_j {
			job.Status = modules.Success_j
		}

		if job.Status == modules.Canceled_j {
			job.Result = "Job canceled"
		}

		//if job.Status == modules.Success_j {
		c := modules.Cluster{Uuid: clusterId}
		_, err = c.Delete()
		if err != nil {
			log.Err(err).Interface("cluster", cluster).Msg("modules delete cluster error")
		}
		log.Debug().Str("clusterUuid", clusterId).Msg("modules delete cluster success")
		//}

		//// rollBACK
		//if job.Status == modules.Failed_j {
		//	cluster.Status = modules.Running_c
		//	err = modules.UpdateByUUID(cluster_old, job.ClusterId)
		//	if err != nil {
		//		log.Error().Err(err).Interface("cluster", cluster).Msg("rollBACK cluster error")
		//	}
		//	log.Debug().Str("cluster uuid", cluster.Uuid).Msg("rollBACK cluster success")
		//}

		job.EndTime = time.Now()
		_, err = job.UpdateByUuid(job.Uuid)
		if err != nil {
			log.Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}
		deleteJob(job)
		if job.Status == modules.Success_j {
			log.Debug().Interface("job", job).Msg("job run success")
		} else {
			log.Warn().Interface("job", job).Msg("job run failed")
		}
	}()

	return job, nil
}

func install(fc *modules.Cluster, values []byte) error {

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
	fc.Spec = result.Config

	return nil
}

func upgrade(fc *modules.Cluster, clusterArgs *ClusterArgs) error {

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
func uninstall(fc *modules.Cluster) error {

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

func Run(j Job) (*modules.Job, error) {
	var clusterArgs *ClusterArgs
	//if ok := service.IsExited(clusterArgs.Name, clusterArgs.Namespace); !ok {
	//	return nil, errors.New("cluster is exited")
	//}
	job := modules.NewJob("ClusterInstall", "")

	//  save job to modules
	_, err := job.Insert()
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("create a new job of ClusterInstall")

	go func() {
		//create a cluster use parameter
		cluster := modules.NewCluster(clusterArgs.Name, clusterArgs.Namespace, clusterArgs.ChartName, clusterArgs.ChartVersion)
		job.ClusterId = cluster.Uuid

		err := install(cluster, clusterArgs.Data)
		if err != nil {
			job.Result = err.Error()
			job.Status = modules.Failed_j
			log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("install cluster error")
		} else {
			job.Result = "install success"
			job.Status = modules.Running_j
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("install cluster success")
		}

		// todo job start status stop timeout
		for job.Status == modules.Running_j {
			clusterStatusOk, err := service.CheckClusterStatus(clusterArgs.Name, clusterArgs.Namespace)
			if err != nil {
				job.Status = modules.Failed_j
				break
			}
			if clusterStatusOk {
				job.Status = modules.Success_j
				break
			}
			time.Sleep(time.Second)
		}

		//todo save cluster to modules
		if job.Status == modules.Success_j {
			_, err = cluster.Insert()
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("save cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("create cluster success")
		}

		job.EndTime = time.Now()
		_, err = job.UpdateByUuid(job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}

		log.Debug().Interface("job", job).Msg("job run success")
	}()

	return job, nil
}
