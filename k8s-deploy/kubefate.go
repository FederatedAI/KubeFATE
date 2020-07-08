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
package main

import (
	"fmt"
	"os"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/cli"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"
)

func main() {
	if err := config.InitConfig(); err != nil {
		fmt.Printf("Unable to read in configuration: %s\n", err)
		os.Exit(1)
	}

	logging.InitLog()

	err := cli.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
