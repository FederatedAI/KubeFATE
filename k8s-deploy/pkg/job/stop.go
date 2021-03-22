/*
 * Copyright 2019-2021 VMware, Inc.
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
	"fmt"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/rs/zerolog/log"
)

// Stop Job Stop
func Stop(jobID string) error {

	j := modules.Job{Uuid: jobID}
	job, err := j.Get()
	if err != nil {
		return fmt.Errorf("get job by jobID error, please check job ID, %s", err.Error())
	}

	// check Method
	if job.Method != "ClusterInstall" {
		return fmt.Errorf("%s type jobs do not support stop", job.Method)
	}

	// check Status
	if job.Status != modules.JobStatusRunning {
		return fmt.Errorf("job status is %s, not support stop", job.Status)
	}

	// set status
	dbErr := job.SetStatus(modules.JobStatusStopping)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("job.SetStatus error")
		return dbErr
	}

	log.Debug().Str("jobID", jobID).Msg("JobStop")

	return nil

}
