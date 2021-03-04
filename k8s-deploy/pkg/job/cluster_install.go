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
	"fmt"
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/rs/zerolog/log"
)

// ClusterInstall Cluster Install New, Create and run job
func ClusterInstall(clusterArgs *modules.ClusterArgs, creator string) (*modules.Job, error) {

	// Preconditions for running job
	// - check cluster args
	// - cluster.IsExisted
	if err := preconditionsReady(clusterArgs, creator); err != nil {
		return nil, err
	}
	log.Debug().Msg("preconditions ready")

	// init job
	job, err := initJob(clusterArgs, "ClusterInstall", creator)
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("init Job success")

	go clusterInstallRun(job)

	return job, nil
}

// PreconditionsReady PreconditionsReady
func preconditionsReady(clusterArgs *modules.ClusterArgs, creator string) error {

	if clusterArgs.ChartVersion == "" {
		return fmt.Errorf("chartVersion cannot be empty")
	}

	if clusterArgs.ChartName == "" {
		return fmt.Errorf("chartVersion cannot be empty")
	}

	cluster := new(modules.Cluster)
	if ok := cluster.IsExisted(clusterArgs.Name, clusterArgs.Namespace); ok {
		return fmt.Errorf("Name=%s & Namespace=%s cluster is existed", clusterArgs.Name, clusterArgs.Namespace)
	}
	return nil
}

func initJob(clusterArgs *modules.ClusterArgs, method, creator string) (*modules.Job, error) {
	//create a job
	job := modules.NewJob(clusterArgs, "ClusterInstall", creator, "")
	//  save job to db
	_, err := job.Insert()
	if err != nil {
		log.Error().Err(err).Interface("job", job).Msg("save job error")
		return nil, err
	}
	return job, nil
}

func clusterInstallRun(job *modules.Job) {

	log.Debug().Str("jobID", job.Uuid).Msg("job Running")
	// update status Running
	err := updateJobStatusToRunning(job)
	// update Event
	if err != nil {
		addJobEvent(job, "update job status to Running, "+err.Error())
		log.Error().Str("jobID", job.Uuid).Err(err).Msg("update job.status to Running")
		return
	}
	addJobEvent(job, "update job status to Running")
	log.Debug().Str("jobID", job.Uuid).Msg("update job.status to Running")

	// create cluster to db
	cluster, err := createCluster(job)

	// update Event
	if err != nil {
		addJobEvent(job, "create cluster error, "+err.Error())
		log.Error().Str("jobID", job.Uuid).Err(err).Msg("create cluster in DB")
		return
	}
	addJobEvent(job, "create cluster in DB Success")
	log.Debug().Str("jobID", job.Uuid).Msg("create Cluster in DB success")

	// helm install a cluster
	err = helmInstall(job, cluster)

	// update Event
	if err != nil {
		clean(job, cluster)
		addJobEvent(job, "helm install error, "+err.Error())
		log.Error().Str("jobID", job.Uuid).Err(err).Msg("Helm install")
		return
	}
	addJobEvent(job, "helm install success")
	log.Debug().Str("jobID", job.Uuid).Msg("Helm install success")

	i := 0
	addJobEvent(job, fmt.Sprintf("checkout cluster status [%d]", i))
	for job.IsRunning() {
		i++
		updateLastJobEvent(job, fmt.Sprintf("checkout cluster status [%d]", i))

		e := &modules.Job{Uuid: job.Uuid}
		j, err := e.Get()
		if err != nil {
			log.Error().Err(err).Msg("get job error")
			return
		}
		job = &j

		// timeOut
		if job.TimeOut() {
			addJobEvent(job, "checkout cluster status timeOut!")
			log.Debug().Str("jobID", job.Uuid).Msg("checkout cluster status timeOut!")
			dbErr := job.SetStatus(modules.JobStatusFailed)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job setStatus error")
			}
			clean(job, cluster)
			break
		}

		// check stop
		if job.IsStop() {
			clean(job, cluster)

			addJobEvent(job, "job stoped")

			dbErr := job.SetStatus(modules.JobStatusCanceled)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetStatus error")
				return
			}

			log.Debug().Str("jobID", job.Uuid).Msg("job stoped")
			break
		}

		// check cluster or job delete
		if stopJob(job, cluster) {
			clean(job, cluster)

			addJobEvent(job, "cluster delete")
			dbErr := job.SetStatus(modules.JobStatusFailed)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetStatus error")
				return
			}
			log.Debug().Str("jobID", job.Uuid).Msg("cluster delete")
			break
		}
		log.Debug().Str("jobID", job.Uuid).Msg("check job Status")
		//check pod status Running , update subJob status
		if checkStatus(job, cluster) {
			addJobEvent(job, "job run Success")

			log.Debug().Str("jobID", job.Uuid).Msg("job run Success")
			{
				cluster.Status = modules.ClusterStatusRunning
				cluster.Revision++
				_, err = cluster.UpdateByUuid(job.ClusterId)
				if err != nil {
					log.Error().Err(err).Interface("cluster", cluster).Msg("update cluster error")
				}
				log.Debug().Str("chart_version", cluster.ChartVersion).Str("status", cluster.Status.String()).Str("cluster uuid", cluster.Uuid).Msg("update cluster success")
			}

			{
				job.Status = modules.JobStatusSuccess
				job.EndTime = time.Now()
				_, err = job.UpdateByUuid(job.Uuid)
				if err != nil {
					log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
				}
			}

			break
		}

		time.Sleep(1 * time.Second)
	}

	if job.Status == modules.JobStatusSuccess {
		log.Debug().Interface("job", job).Msg("job run success")
	} else {
		log.Warn().Interface("job", job).Msg("job run failed")
	}

}

func createCluster(job *modules.Job) (*modules.Cluster, error) {
	// create a cluster
	cluster, err := modules.NewCluster(job.Metadata.Name, job.Metadata.Namespace, job.Metadata.ChartName, job.Metadata.ChartVersion, string(job.Metadata.Data))
	if err != nil {
		log.Error().Err(err).Interface("clusterArgs", job.Metadata).Msg("NewCluster")
		return nil, err
	}

	// Save cluster to database
	_, err = cluster.Insert()
	if err != nil {
		log.Error().Err(err).Interface("cluster", cluster).Msg("save cluster error")
		return nil, err
	}

	dbErr := cluster.SetStatus(modules.ClusterStatusCreating)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("cluster setStatus error")
	}

	// update ClusterId of job
	job.ClusterId = cluster.Uuid
	_, dbErr = job.UpdateByUuid(job.Uuid)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("cluster setStatus error")
	}

	log.Info().Str("cluster uuid", cluster.Uuid).Msg("save cluster Success")
	return cluster, nil
}

func updateJobStatusToRunning(job *modules.Job) error {
	e := &modules.Job{Uuid: job.Uuid}
	j, dbErr := e.Get()
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("get job error")
		return dbErr
	}
	job = &j
	if job.Status != modules.JobStatusPending {
		return errors.New("job.Status error")
	}
	dbErr = job.SetStatus(modules.JobStatusRunning)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("job setStatus error")
		return dbErr
	}
	dbErr = job.SetEvent("job start Running")
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("job setResult error")
		return dbErr
	}
	return nil
}

func helmInstall(job *modules.Job, cluster *modules.Cluster) error {

	clusterCover(job, cluster)

	err := cluster.HelmInstall()
	if err != nil {
		log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("helm install cluster error")

		dbErr := job.SetStatus(modules.JobStatusFailed)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("job.SetStatus error")
		}
		return err
	}

	dbErr := job.SetStatus(modules.JobStatusRunning)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("job.SetStatus error")
	}
	return nil

}

func clean(job *modules.Job, cluster *modules.Cluster) {

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

func addJobEvent(job *modules.Job, Event string) {

	job.Events = append(job.Events, Event)

	dbErr := job.SetEvents(job.Events)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("job.SetStatus error")
	}
	return
}
func updateLastJobEvent(job *modules.Job, Event string) {

	job.Events[len(job.Events)-1] = Event

	dbErr := job.SetEvents(job.Events)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("job.SetStatus error")
	}
	return
}
func checkStatus(job *modules.Job, cluster *modules.Cluster) bool {

	// update subJobs
	ClusterStatus, err := service.GetClusterPodStatus(cluster.Name, cluster.NameSpace)
	if err != nil {
		log.Error().Err(err).Msg("GetClusterPodStatus error")
	}

	log.Debug().Interface("ClusterStatus", ClusterStatus).Msg("GetClusterPodStatus()")
	subJobs := generateSubJobs(job, ClusterStatus)

	dbErr := job.SetSubJobs(subJobs)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("job setSubJobs error")
	}

	if service.CheckClusterStatus(ClusterStatus) {
		dbErr := job.SetStatus(modules.JobStatusSuccess)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("job setStatus error")
		}
		return true
	}
	return false
}

func clusterCover(job *modules.Job, cluster *modules.Cluster) {
	// Override existing installation
	if job.Metadata.Cover {
		log.Info().Msg("overwrite current installation")
		err := cluster.HelmDelete()
		if err != nil {
			log.Error().Err(err).Msg("helm clean error")
		}
		addJobEvent(job, "overwrite current installation")
		log.Info().Str("name", cluster.Name).Str("namespace", cluster.NameSpace).Msg("HelmClean success")
	}
}
