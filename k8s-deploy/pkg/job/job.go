/*
 * Copyright 2019-2022 VMware, Inc.
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
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/hashicorp/go-version"
	"github.com/rs/zerolog/log"
)

const (
	fateChartName               = "fate"
	fateUpgradeManagerChartName = "fate-upgrade-manager"
	fateUpgradeJobName          = "fate-mysql-upgrade-job"
)

// validateFateVersion helps make sure that the user set the right helm chart version and
// the image version is equal to the chart version. For the versions not in the keys
// of chartToImageVersionMap, we just skip the validation because this KubeFATE service
// should also support some future versions.
func validateFateVersion(chartVersion string, imageVersion string) error {
	chartToImageVersionMap := map[string]string{
		"v1.7.0": "1.7.0-release",
		"v1.7.1": "1.7.1-release",
		"v1.7.2": "1.7.2-release",
		"v1.8.0": "1.8.0-release",
		"v1.9.0": "1.9.0-release",
	}
	if expectedImageVersion, ok := chartToImageVersionMap[chartVersion]; ok {
		if expectedImageVersion == imageVersion {
			return nil
		}
		log.Error().Msgf("the chart version is %s but the image version is %s", chartVersion, imageVersion)
		return errors.New("the image tag is not consistent with the chart version, which is unsupported")
	}
	return nil
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

func generateSubJobs_b(job *modules.Job, ClusterStatus map[string]string) modules.SubJobs {

	subJobs := make(modules.SubJobs)
	if job.SubJobs != nil {
		subJobs = job.SubJobs
	}

	for k, v := range ClusterStatus {
		var subJobStatus string
		if v == "Running" {
			subJobStatus = "Success"
		} else if v == "Failed" || v == "Unknown" || v == "Pending" {
			subJobStatus = v
		} else {
			subJobStatus = "Running"
		}

		var subJob modules.SubJob
		if _, ok := subJobs[k]; !ok {
			subJob = modules.SubJob{
				ModuleName:    k,
				Status:        subJobStatus,
				ModulesStatus: v,
				StartTime:     job.StartTime,
			}
		} else {
			subJob = subJobs[k]
			subJob.Status = subJobStatus
			subJob.ModulesStatus = v
		}

		if subJobStatus == "Success" && subJob.EndTime.IsZero() {
			subJob.EndTime = time.Now()
		}

		subJobs[k] = subJob
		log.Debug().Interface("subJob", subJob).Msg("generate SubJobs")
	}

	job.SubJobs = subJobs
	return subJobs
}

func generateSubJobs(job *modules.Job, ClusterDeployStatus map[string]string) modules.SubJobs {

	subJobs := make(modules.SubJobs)
	if job.SubJobs != nil {
		subJobs = job.SubJobs
	}

	for k, v := range ClusterDeployStatus {
		var subJobStatus string = "Running"
		if service.CheckStatus(v) {
			subJobStatus = "Success"
		}

		var subJob modules.SubJob
		if _, ok := subJobs[k]; !ok {
			subJob = modules.SubJob{
				ModuleName:    k,
				Status:        subJobStatus,
				ModulesStatus: v,
				StartTime:     job.StartTime,
			}
		} else {
			subJob = subJobs[k]
			subJob.Status = subJobStatus
			subJob.ModulesStatus = v
		}

		if subJobStatus == "Success" && subJob.EndTime.IsZero() {
			subJob.EndTime = time.Now()
		}

		subJobs[k] = subJob
		log.Debug().Interface("subJob", subJob).Msg("generate SubJobs")
	}

	job.SubJobs = subJobs
	return subJobs
}

func ClusterUpdate(clusterArgs *modules.ClusterArgs, creator string) (*modules.Job, error) {
	// Check whether the cluster exists
	c := new(modules.Cluster)
	if ok := c.IsExisted(clusterArgs.Name, clusterArgs.Namespace); !ok {
		return nil, fmt.Errorf("name=%s Cluster is not existed", clusterArgs.Name)
	}

	c = &modules.Cluster{Name: clusterArgs.Name, NameSpace: clusterArgs.Namespace}
	cluster, err := c.Get()
	if err != nil {
		log.Error().Err(err).Interface("clusterArgs", clusterArgs).Msg("Find Cluster by clusterArgs error")
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

	if cluster.ChartName != clusterNew.ChartName {
		return nil, fmt.Errorf("doesn't support upgrade between different charts")
	}

	if reflect.DeepEqual(specOld, specNew) &&
		cluster.ChartName == clusterArgs.ChartName &&
		cluster.ChartVersion == clusterArgs.ChartVersion {
		return nil, fmt.Errorf("the configuration file did not change")
	}

	oldVersion, err := version.NewVersion(strings.ReplaceAll(cluster.ChartVersion, "v", ""))
	newVersion, err := version.NewVersion(strings.ReplaceAll(clusterArgs.ChartVersion, "v", ""))
	// Comparison example. There is also GreaterThan, Equal, and just
	// a simple Compare that returns an int allowing easy >=, <=, etc.
	if newVersion.LessThan(oldVersion) {
		return nil, fmt.Errorf("using kubefate to downgrading a cluster is not supported")
	}
	// FmlFrameWorkNameToUmStuffMap is a map, whose keys are the fml framework names, such as fate.
	// Its values are the Information of the corresponding upgrade manager's chart.
	FmlFrameWorkNameToUmInfoMap := modules.MapStringInterface{
		fateChartName: modules.MapStringInterface{
			"chartName":             fateUpgradeManagerChartName,
			"chartVersion":          "v1.0.0",
			"specConstructFunction": ConstructFumSpec,
		},
	}
	var upgradeManagerChartName string
	if newVersion.GreaterThan(oldVersion) {
		switch clusterNew.ChartName {
		case fateChartName:
			err = validateFateVersion(c.ChartVersion, specNew["imageTag"].(string))
			if err != nil {
				return nil, err
			}
			upgradeManagerChartName = fateUpgradeManagerChartName
		default:
			log.Warn().Msgf("the chart %s doesn't have an upgrade manager", clusterNew.ChartName)
		}
	}
	job := modules.NewJob(clusterArgs, "ClusterUpdate", creator, cluster.Uuid)
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
			log.Error().Err(dbErr).Msg("Cluster.SetStatus error")
		}

		// If the chart version changed, we will get a not nil chart name here
		// We will implicitly install a new cluster for the upgrade manager, and delete it after it finishes its job
		if upgradeManagerChartName != "" {
			umInfo := FmlFrameWorkNameToUmInfoMap[cluster.ChartName].(modules.MapStringInterface)
			// Use reflection to run the special logic of the upgrade manager for each fml framework
			umSpecFuncName := reflect.ValueOf(umInfo["specConstructFunction"])
			params := []reflect.Value{reflect.ValueOf(specOld), reflect.ValueOf(specNew)}
			res := umSpecFuncName.Call(params)
			umSpec := res[0].Interface().(modules.MapStringInterface)
			umCluster := modules.Cluster{
				Name:         upgradeManagerChartName,
				NameSpace:    cluster.NameSpace,
				ChartName:    upgradeManagerChartName,
				ChartVersion: umInfo["chartVersion"].(string),
				Spec:         umSpec,
			}
			err := umCluster.HelmInstall()
			if err != nil {
				log.Error().Err(err).Msgf("failed to install the upgrade manager's helm chart for cluster %s", cluster.ChartName)
				dbErr := job.SetStatus(modules.JobStatusFailed)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetStatus error")
				}
			}
			var cycle int
			interval := 30
			for cycle = 0; cycle < 20; cycle++ {
				jobDone, err := service.CheckJobReadiness(umCluster.NameSpace, fateUpgradeJobName)
				if err != nil {
					log.Error().Err(err).Msg("failed to check the upgrade job status")
				}
				if jobDone {
					log.Info().Msgf("the upgrade manager finished its job in %d seconds", interval*cycle)
					break
				}
				log.Info().Msgf("the upgrade manager's job is not done yet, will recheck in %d seconds", interval)
				time.Sleep(time.Second * time.Duration(interval))
			}
			if cycle == 20 {
				errMsg := fmt.Sprintf("the upgrade manager cannot finish the job in %d seconds", 30*cycle)
				err := errors.New(errMsg)
				log.Error().Err(err)
				// we will do helm delete to the upgrade manager if this time out is triggered, then
				// we just need to set the job to failed, and later there is a logic will handle
				// the rollback of the FML pods
				dbErr := job.SetStatus(modules.JobStatusFailed)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetStatus error")
				}
			}
			if !clusterArgs.KeepUpgradeJob {
				err = umCluster.HelmDelete()
				if err != nil {
					log.Error().Err(err).Msg("failed to delete the upgrade manager cluster, need a person to investigate why")
				}
			}
		}

		dbErr = cluster.SetValues(valuesNew)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("Cluster.SetSpec error")
		}
		dbErr = cluster.SetSpec(specNew)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("Cluster.SetSpec error")
		}

		// HelmUpgrade

		//The Chart version does not change and update is used
		//Upgrade corresponding to Helm
		cluster.ChartName = clusterArgs.ChartName
		cluster.ChartVersion = clusterArgs.ChartVersion
		err = cluster.HelmUpgrade()
		cluster.HelmRevision += 1

		_, dbErr = cluster.UpdateByUuid(job.ClusterId)
		if err != nil {
			log.Error().Err(dbErr).Interface("cluster", cluster).Msg("Update Cluster error")
		}

		if err != nil {
			log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("Helm upgrade Cluster error")

			dbErr := job.SetState(err.Error())
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetResult error")
			}
			dbErr = job.SetStatus(modules.JobStatusFailed)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetStatus error")
			}
		} else {
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("helm upgrade Cluster Success")

			dbErr := job.SetState("Cluster update Success")
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
				dbErr := job.SetState("Checkout Cluster status timeOut!")
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
			ClusterStatus, err := service.GetClusterDeployStatus(clusterArgs.Name, clusterArgs.Namespace)
			if err != nil {
				log.Error().Err(err).Msg("GetClusterDeployStatus error")
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
			dbErr := job.SetState("Job canceled")
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetResult error")
			}
		}

		// save cluster to modules
		if job.Status == modules.JobStatusSuccess {
			cluster.Status = modules.ClusterStatusRunning
			cluster.Revision++
			_, err = cluster.UpdateByUuid(job.ClusterId)
			if err != nil {
				log.Error().Err(err).Interface("cluster", cluster).Msg("Update Cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("Update Cluster Success")
		}

		// rollBACK
		if job.Status != modules.JobStatusSuccess && job.Status != modules.JobStatusCanceled {
			dbErr = cluster.SetValues(valuesOld)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("Cluster.SetSpec error")
			}
			dbErr = cluster.SetSpec(specOld)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("Cluster.SetSpec error")
			}
			dbErr = cluster.SetStatus(modules.ClusterStatusRollback)
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("Cluster.SetStatus error")
			}

			//The Chart version does not change and update is used
			//Upgrade corresponding to Helm
			err = cluster.HelmRollback()
			cluster.HelmRevision -= 1

			if err != nil {
				log.Error().Err(err).Str("ClusterId", cluster.Uuid).Msg("Helm upgrade Cluster error")

				dbErr := job.SetState(err.Error())
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetResult error")
				}
				dbErr = job.SetStatus(modules.JobStatusFailed)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job.SetStatus error")
				}
			} else {
				log.Debug().Str("ClusterId", cluster.Uuid).Msg("Helm upgrade Cluster Success")

				dbErr := job.SetState("Cluster run rollback Success")
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
					dbErr := job.SetState("Checkout Cluster status timeOut!")
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
				ClusterStatus, err := service.GetClusterDeployStatus(clusterArgs.Name, clusterArgs.Namespace)
				if err != nil {
					log.Error().Err(err).Msg("GetClusterDeployStatus error")
				}

				log.Debug().Interface("ClusterStatus", ClusterStatus).Msg("GetClusterDeployStatus()")

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
				log.Error().Err(err).Interface("cluster", cluster).Msg("RollBACK Cluster error")
			}
			log.Debug().Str("cluster uuid", cluster.Uuid).Msg("RollBACK Cluster Success")
		}

		job.EndTime = time.Now()
		_, err = job.UpdateByUuid(job.Uuid)
		if err != nil {
			log.Error().Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}

		if job.Status == modules.JobStatusSuccess {
			log.Debug().Interface("job", job).Msg("job run Success")
		} else {
			log.Warn().Interface("job", job).Msg("job run failed")
		}
	}()

	return job, nil
}

func ClusterDelete(clusterId string, creator string) (*modules.Job, error) {
	if clusterId == "" {
		return nil, fmt.Errorf("clusterID cannot be empty")
	}

	c := modules.Cluster{Uuid: clusterId}
	cluster, err := c.Get()
	if err != nil {
		log.Error().Err(err).Interface("clusterID", clusterId).Msg("Find Cluster by clusterId error")
		return nil, err
	}

	job := modules.NewJob(nil, "ClusterDelete", creator, clusterId)
	// save job to modules
	_, err = job.Insert()
	if err != nil {
		log.Err(err).Interface("job", job).Msg("Save job error")
		return nil, err
	}

	log.Info().Str("jobId", job.Uuid).Msg("Create a new job of ClusterDelete")

	go func() {
		dbErr := job.SetStatus(modules.JobStatusRunning)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("job.SetStatus error")
		}
		dbErr = cluster.SetStatus(modules.ClusterStatusDeleting)
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("Cluster.SetStatus error")
		}

		err = cluster.HelmDelete()
		if err != nil {
			dbErr := job.SetState(err.Error())
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetResult error")
			}
			job.Status = modules.JobStatusFailed
			log.Err(err).Str("ClusterId", cluster.Uuid).Msg("Helm delete Cluster error")
		} else {
			dbErr := job.SetState("uninstall Success")
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetResult error")
			}
			job.Status = modules.JobStatusRunning
			log.Debug().Str("ClusterId", cluster.Uuid).Msg("Helm delete Cluster Success")
		}

		if job.Status == modules.JobStatusRunning {
			job.Status = modules.JobStatusSuccess
		}

		if job.Status == modules.JobStatusCanceled {
			dbErr := job.SetState("Job canceled")
			if dbErr != nil {
				log.Error().Err(dbErr).Msg("job.SetResult error")
			}
		}

		//if job.Status == modules.JobStatusSuccess {
		c := modules.Cluster{Uuid: clusterId}
		_, err = c.Delete()
		if err != nil {
			log.Err(err).Interface("cluster", cluster).Msg("modules delete Cluster error")
		}
		log.Debug().Str("clusterUuid", clusterId).Msg("modules delete Cluster Success")

		job.EndTime = time.Now()
		_, err = job.UpdateByUuid(job.Uuid)
		if err != nil {
			log.Err(err).Str("jobId", job.Uuid).Msg("update job By Uuid error")
		}
		if job.Status == modules.JobStatusSuccess {
			log.Debug().Interface("job", job).Msg("job run Success")
		} else {
			log.Warn().Interface("job", job).Msg("job run failed")
		}
	}()

	return job, nil
}
