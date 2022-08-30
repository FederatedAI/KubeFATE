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

package job

import (
	"errors"
	"fmt"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/hashicorp/go-version"
	"github.com/rs/zerolog/log"
	"reflect"
	"strings"
	"time"
)

const (
	fumChartName = "fate-upgrade-manager"
	fumJobName   = "fate-mysql-upgrade-job"
)

type FateUpgradeManager struct {
	namespace string
	UpgradeManager
}

func (fum *FateUpgradeManager) validate(specOld, specNew modules.MapStringInterface) error {
	oldChartName := specOld["chartName"].(string)
	newChartName := specNew["chartName"].(string)
	oldVersion := specOld["chartVersion"].(string)
	newVersion := specNew["chartVersion"].(string)

	// Chart must be the same fml framework
	if oldChartName != newChartName {
		return errors.New("doesn't support upgrade between different charts")
	}

	// To upgrade, chart must have some difference
	if reflect.DeepEqual(specOld, specNew) && oldVersion == newVersion {
		return errors.New("the configuration file did not change")
	}

	// Do not support downgrade
	newVerFormatted, _ := version.NewVersion(newVersion)
	oldVerFormatted, _ := version.NewVersion(oldVersion)
	if newVerFormatted.LessThan(oldVerFormatted) {
		return errors.New("using kubefate to downgrade a cluster is not supported yet")
	}

	// KubeFATE cannot support rolling upgrade for FATE version <= 1.7.1
	if newVersion != oldVersion {
		ver171, _ := version.NewVersion("1.7.1")
		if oldVerFormatted.LessThanOrEqual(ver171) {
			return errors.New("upgrade from FATE version <= 1.7.1 is not supported by KubeFATE")
		}
	}
	log.Info().Msg("version validation for FATE cluster yaml passed")
	return nil
}

func (fum *FateUpgradeManager) getCluster(specOld, specNew modules.MapStringInterface) modules.Cluster {
	fumCluster := modules.Cluster{
		Name:         fumChartName,
		NameSpace:    fum.namespace,
		ChartName:    fumChartName,
		ChartVersion: "v1.0.0",
		Spec:         constructFumSpec(specOld, specNew),
	}
	return fumCluster
}

func (fum *FateUpgradeManager) waitFinish(interval, round int) bool {
	var cycle int
	for cycle = 0; cycle < round; cycle++ {
		jobDone, err := service.CheckJobReadiness(fum.namespace, fumJobName)
		if err != nil {
			log.Error().Err(err).Msg("failed to check the upgrade job status")
		}
		if jobDone {
			log.Info().Msgf("the upgrade manager finished its job in %d seconds", interval*cycle)
			return true
		}
		log.Info().Msgf("the upgrade manager's job is not done yet, will recheck in %d seconds", interval)
		time.Sleep(time.Second * time.Duration(interval))
	}
	errMsg := fmt.Sprintf("the upgrade manager cannot finish the job in %d seconds", 30*cycle)
	log.Error().Msg(errMsg)
	return false
}

func getMysqlCredFromSpec(clusterSpec modules.MapStringInterface) (username, password string) {
	mysqlSpec := clusterSpec["mysql"].(map[string]interface{})
	if mysqlSpec["user"] == nil {
		username = "fate"
	} else {
		username = mysqlSpec["user"].(string)
	}
	if mysqlSpec["password"] == nil {
		password = "fate_dev"
	} else {
		password = mysqlSpec["password"].(string)
	}
	return
}

func constructFumSpec(oldSpec, newSpec modules.MapStringInterface) (fumSpec modules.MapStringInterface) {
	oldVersion := strings.ReplaceAll(oldSpec["chartVersion"].(string), "v", "")
	newVersion := strings.ReplaceAll(newSpec["chartVersion"].(string), "v", "")
	mysqlUsername, mysqlPassword := getMysqlCredFromSpec(newSpec)
	res := modules.MapStringInterface{
		"username": mysqlUsername,
		"password": mysqlPassword,
		"start":    oldVersion,
		"target":   newVersion,
	}
	return res
}
