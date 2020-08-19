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
	"reflect"
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

	// Check whether the cluster exists
	cluster := new(modules.Cluster)
	if ok := cluster.IsExisted(clusterArgs.Name, clusterArgs.Namespace); ok {
		return nil, fmt.Errorf("name=%s cluster is exited", clusterArgs.Name)
	}

	// create a cluster
	cluster, err := modules.NewCluster(clusterArgs.Name, clusterArgs.Namespace, clusterArgs.ChartName, clusterArgs.ChartVersion, string(clusterArgs.Data))
	if err != nil {
		log.Error().Err(err).Interface("clusterArgs", clusterArgs).Msg("NewCluster")
		return nil, err
	}

	// Save cluster to database
	_, err = cluster.Insert()
	if err != nil {
		log.Error().Err(err).Interface("cluster", cluster).Msg("save cluster error")
		return nil, err
	}
	log.Info().Str("cluster uuid", cluster.Uuid).Msg("save cluster success")

	//create a job
	job := modules.NewJob("ClusterInstall", creator, cluster.Uuid)
	//  save job to db
	_, err = job.Insert()
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("create a new job of ClusterInstall")

	//do job
	go func() {
		dbErr := job.SetStatus(modules.JobStatusRunning)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("job.SetStatus error")
		}

		dbErr = cluster.SetStatus(modules.ClusterStatusCreating)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("cluster.SetStatus error")
		}

		// Override existing installation
		if clusterArgs.Cover {
			log.Info().Msg("Overwrite current installation")
			err = cluster.HelmDelete()
			if err != nil {
				log.Error().Err(err).Msg("helmClean error")
			}
			log.Info().Str("name", cluster.Name).Str("namespace", cluster.NameSpace).Msg("HelmClean success")
		}

		err = cluster.HelmInstall()
		if err != nil {
			log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("helm install cluster error")

			dbErr := job.SetResult(err.Error())
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetResult error")
			}
			dbErr = job.SetStatus(modules.JobStatusFailed)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetStatus error")
			}
		} else {
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("helm install cluster success")

			dbErr := job.SetResult("Cluster install success")
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetResult error")
			}
			dbErr = job.SetStatus(modules.JobStatusRunning)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetStatus error")
			}
		}

		//
		for job.Status == modules.JobStatusRunning {

			if !cluster.IsExisted(cluster.Name, cluster.NameSpace) {
				dbErr := job.SetResult("Cluster has been deleted!")
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetResult error")
				}
				dbErr = job.SetStatus(modules.JobStatusCanceled)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetStatus error")
				}
				continue
			}

			if stopJob(job, cluster) {
				continue
			}

			if job.TimeOut() {
				dbErr := job.SetResult("Checkout cluster status timeOut!")
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetResult error")
				}
				dbErr = job.SetStatus(modules.JobStatusFailed)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetStatus error")
				}
				break
			}

			// update subJobs
			ClusterStatus, err := service.GetClusterStatus(clusterArgs.Name, clusterArgs.Namespace)
			if err != nil {
				log.Error().Err(err).Msg("GetClusterStatus error")
			}

			log.Debug().Interface("ClusterStatus", ClusterStatus).Msg("GetClusterStatus()")
			subJobs := generateSubJobs(job, ClusterStatus)

			dbErr = job.SetSubJobs(subJobs)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetSubJobs error")
			}

			if service.CheckClusterStatus(ClusterStatus) {
				dbErr := job.SetStatus(modules.JobStatusSuccess)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetStatus error")
				}
				break
			}
			time.Sleep(1 * time.Second)
		}

		if job.Status == modules.JobStatusCanceled {
			job.Result = "Job canceled"
		}

		// save cluster to modules
		if job.Status == modules.JobStatusSuccess {
			cluster.Status = modules.ClusterStatusRunning
			cluster.Revision += 1
			_, err = cluster.UpdateByUuid(job.ClusterId)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
			}
			log.Debug().Str("chart_version", cluster.ChartVersion).Str("status", cluster.Status.String()).Str("cluster uuid", cluster.Uuid).Msg("update cluster success")
		}

		// rollBACK
		if job.Status != modules.JobStatusSuccess && job.Status != modules.JobStatusCanceled {
			_, err := cluster.Delete()
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("DB delete cluster error")
			}
			err = cluster.HelmDelete()
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("helm delete cluster error")
			}

			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("install cluster rollBACK success")
		}

		job.EndTime = time.Now()
		_, err = job.UpdateByUuid(job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}

		if job.Status == modules.JobStatusSuccess {
			log.Debug().Interface("job", job).Msg("job run success")
		} else {
			log.Warn().Interface("job", job).Msg("job run failed")
		}

	}()

	return job, nil
}

func stopJob(job *modules.Job, cluster *modules.Cluster) bool {
	if !cluster.IsExisted(cluster.Name, cluster.NameSpace) {
		return true
	}

	if !job.IsExisted(job.Uuid) {
		return true
	}

	return false
}

func generateSubJobs(job *modules.Job, ClusterStatus map[string]string) modules.SubJobs {

	subJobs := make(modules.SubJobs)
	if job.SubJobs != nil {
		subJobs = job.SubJobs
	}
	log.Debug().Interface("subJobs", subJobs).Msg("subJobs=job.SubJobs")
	for k, v := range ClusterStatus {
		var subJobStatus string
		if v == "Running" {
			subJobStatus = "Success"
		} else {
			subJobStatus = "Running"
		}

		var subJob modules.SubJob
		if _, ok := subJobs[k]; !ok {
			subJob = modules.SubJob{
				ModuleName:       k,
				Status:           subJobStatus,
				ModulesPodStatus: v,
				StartTime:        job.StartTime,
			}
		} else {
			subJob = subJobs[k]
			subJob.Status = subJobStatus
			subJob.ModulesPodStatus = v
		}

		if subJobStatus == "Success" && subJob.EndTime.IsZero() {
			subJob.EndTime = time.Now()
		}

		subJobs[k] = subJob
	}

	job.SubJobs = subJobs
	return subJobs
}

func ClusterUpdate(clusterArgs *ClusterArgs, creator string) (*modules.Job, error) {
	// Check whether the cluster exists
	c := new(modules.Cluster)
	if ok := c.IsExisted(clusterArgs.Name, clusterArgs.Namespace); !ok {
		return nil, fmt.Errorf("name=%s cluster is not exited", clusterArgs.Name)
	}

	c = &modules.Cluster{Name: clusterArgs.Name, NameSpace: clusterArgs.Namespace}
	cluster, err := c.Get()
	if err != nil {
		log.Error().Err(err).Interface("clusterArgs", clusterArgs).Msg("find cluster by clusterArgs error")
		return nil, err
	}

	clusterNew, err := modules.NewCluster(clusterArgs.Name, clusterArgs.Namespace, clusterArgs.ChartName, clusterArgs.ChartVersion, string(clusterArgs.Data))
	if err != nil {
		log.Error().Err(err).Msg("NewCluster")
		return nil, err
	}

	var specOld = cluster.Spec
	var specNew = clusterNew.Spec
	var valuesOld = cluster.Values
	var valuesNew = clusterNew.Values

	if reflect.DeepEqual(specOld, specNew) &&
		cluster.ChartName == clusterArgs.ChartName &&
		cluster.ChartVersion == clusterArgs.ChartVersion {
		return nil, fmt.Errorf("the configuration file did not change")
	}

	job := modules.NewJob("ClusterUpdate", creator, cluster.Uuid)
	//  save job to modules
	_, err = job.Insert()
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("create a new job of ClusterUpdate")

	//do job
	go func() {
		// update job.status/ cluster.status / cluster
		dbErr := job.SetStatus(modules.JobStatusRunning)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("job.SetStatus error")
		}

		dbErr = cluster.SetStatus(modules.ClusterStatusUpdating)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("cluster.SetStatus error")
		}

		dbErr = cluster.SetValues(valuesNew)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("cluster.SetSpec error")
		}
		dbErr = cluster.SetSpec(specNew)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("cluster.SetSpec error")
		}

		// HelmUpgrade

		//The chart version does not change and update is used
		//Upgrade corresponding to Helm
		cluster.ChartName = clusterArgs.ChartName
		cluster.ChartVersion = clusterArgs.ChartVersion
		err = cluster.HelmUpgrade()
		cluster.HelmRevision += 1

		_, dbErr = cluster.UpdateByUuid(job.ClusterId)
		if err != nil {
			log.Error().Err(dbErr).Interface("cluster", cluster).Msg("update cluster error")
		}

		if err != nil {
			log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("helm upgrade cluster error")

			dbErr := job.SetResult(err.Error())
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetResult error")
			}
			dbErr = job.SetStatus(modules.JobStatusFailed)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetStatus error")
			}
		} else {
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("helm upgrade cluster success")

			dbErr := job.SetResult("Cluster update success")
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetResult error")
			}
			dbErr = job.SetStatus(modules.JobStatusRunning)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetStatus error")
			}
		}

		//
		for job.Status == modules.JobStatusRunning {
			if stopJob(job, &cluster) {
				continue
			}

			if job.TimeOut() {
				dbErr := job.SetResult("Checkout cluster status timeOut!")
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetResult error")
				}
				dbErr = job.SetStatus(modules.JobStatusFailed)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetStatus error")
				}
				break
			}

			// update subJobs
			ClusterStatus, err := service.GetClusterStatus(clusterArgs.Name, clusterArgs.Namespace)
			if err != nil {
				log.Error().Err(err).Msg("GetClusterStatus error")
			}

			subJobs := generateSubJobs(job, ClusterStatus)

			dbErr = job.SetSubJobs(subJobs)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetSubJobs error")
			}

			if service.CheckClusterStatus(ClusterStatus) {
				dbErr := job.SetStatus(modules.JobStatusSuccess)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetStatus error")
				}
				break
			}
			time.Sleep(5 * time.Second)
		}

		if job.Status == modules.JobStatusCanceled {
			job.Result = "Job canceled"
		}

		// save cluster to modules
		if job.Status == modules.JobStatusSuccess {
			cluster.Status = modules.ClusterStatusRunning
			cluster.Revision += 1
			_, err = cluster.UpdateByUuid(job.ClusterId)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("update cluster success")
		}

		// rollBACK
		if job.Status != modules.JobStatusSuccess && job.Status != modules.JobStatusCanceled {
			//todo helm rollBack
			dbErr = cluster.SetValues(valuesOld)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("cluster.SetSpec error")
			}
			dbErr = cluster.SetSpec(specOld)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("cluster.SetSpec error")
			}
			dbErr = cluster.SetStatus(modules.ClusterStatusRollback)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("cluster.SetStatus error")
			}

			//The chart version does not change and update is used
			//Upgrade corresponding to Helm
			err = cluster.HelmRollback()
			cluster.HelmRevision -= 1

			if err != nil {
				log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("helm upgrade cluster error")

				dbErr := job.SetResult(job.Result + err.Error())
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetResult error")
				}
				dbErr = job.SetStatus(modules.JobStatusFailed)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetStatus error")
				}
			} else {
				log.Debug().Str("ClusterId", cluster.Uuid).Msg("helm upgrade cluster success")

				dbErr := job.SetResult(job.Result + "Cluster run rollback success")
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetResult error")
				}
				dbErr = job.SetStatus(modules.JobStatusRollback)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetStatus error")
				}
			}

			//
			for job.Status == modules.JobStatusRunning {
				if job.TimeOut() {
					dbErr := job.SetResult("Checkout cluster status timeOut!")
					if dbErr != nil {
						log.Error().Err(dbErr).Msg("job.SetResult error")
					}
					dbErr = job.SetStatus(modules.JobStatusFailed)
					if dbErr != nil {
						log.Error().Err(dbErr).Msg("job.SetStatus error")
					}
					break
				}

				// update subJobs
				ClusterStatus, err := service.GetClusterStatus(clusterArgs.Name, clusterArgs.Namespace)
				if err != nil {
					log.Error().Err(err).Msg("GetClusterStatus error")
				}

				log.Debug().Interface("ClusterStatus", ClusterStatus).Msg("GetClusterStatus()")

				subJobs := generateSubJobs(job, ClusterStatus)

				dbErr = job.SetSubJobs(subJobs)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetSubJobs error")
				}

				if service.CheckClusterStatus(ClusterStatus) {
					dbErr := job.SetStatus(modules.JobStatusSuccess)
					if dbErr != nil {
						log.Error().Err(dbErr).Msg("job.SetStatus error")
					}
					break
				}
				time.Sleep(5 * time.Second)
			}

			_, err = cluster.UpdateByUuid(cluster.Uuid)
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

		if job.Status == modules.JobStatusSuccess {
			log.Debug().Interface("job", job).Msg("job run success")
		} else {
			log.Warn().Interface("job", job).Msg("job run failed")
		}
	}()

	return job, nil
}

func ClusterDelete(clusterId string, creator string) (*modules.Job, error) {
	if clusterId == "" {
		return nil, fmt.Errorf("clusterid cannot be empty")
	}

	c := modules.Cluster{Uuid: clusterId}
	cluster, err := c.Get()
	if err != nil {
		log.Error().Err(err).Interface("clusterId", clusterId).Msg("find cluster by clusterId error")
		return nil, err
	}

	job := modules.NewJob("ClusterDelete", creator, clusterId)
	// save job to modules
	_, err = job.Insert()
	if err != nil {
		log.Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("create a new job of ClusterDelete")

	go func() {
		dbErr := job.SetStatus(modules.JobStatusRunning)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("job.SetStatus error")
		}
		dbErr = cluster.SetStatus(modules.ClusterStatusDeleting)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("cluster.SetStatus error")
		}

		err = cluster.HelmDelete()
		if err != nil {
			job.Result = err.Error()
			job.Status = modules.JobStatusFailed
			log.Err(err).Str("ClusterId", cluster.Uuid).Msg("helm delete cluster error")
		} else {
			job.Result = "uninstall success"
			job.Status = modules.JobStatusRunning
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("helm delete cluster success")
		}

		if job.Status == modules.JobStatusRunning {
			job.Status = modules.JobStatusSuccess
		}

		if job.Status == modules.JobStatusCanceled {
			job.Result = "Job canceled"
		}

		//if job.Status == modules.JobStatusSuccess {
		c := modules.Cluster{Uuid: clusterId}
		_, err = c.Delete()
		if err != nil {
			log.Err(err).Interface("cluster", cluster).Msg("modules delete cluster error")
		}
		log.Debug().Str("clusterUuid", clusterId).Msg("modules delete cluster success")

		job.EndTime = time.Now()
		_, err = job.UpdateByUuid(job.Uuid)
		if err != nil {
			log.Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}
		if job.Status == modules.JobStatusSuccess {
			log.Debug().Interface("job", job).Msg("job run success")
		} else {
			log.Warn().Interface("job", job).Msg("job run failed")
		}
	}()

	return job, nil
}
