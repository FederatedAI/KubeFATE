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
package cli

import (
	"os"
	"testing"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"
)

func TestRun(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	type args struct {
		Args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"help",
			args{[]string{os.Args[0], "help"}},
		},
		{
			"cluster -help",
			args{[]string{os.Args[0], "cluster", "--help"}},
		},
		{
			"cluster list",
			args{[]string{os.Args[0], "cluster", "list"}},
		},
		//{
		//	"get -help",
		//	args{[]string{os.Args[0], "get", "--help"}},
		//},
		//{
		//	"delete -help",
		//	args{[]string{os.Args[0], "delete", "--help"}},
		//},
		//{
		//	"list all",
		//	args{[]string{os.Args[0], "list", "--namespace", ""}},
		//},
		//{
		//	"install",
		//	args{[]string{os.Args[0], "install", "--name", "fate-10000", "--namespace", "fate-10000", "--chart", "E:\\machenlong\\AI\\github\\owlet42\\KubeFATE\\k8s-deploy\\fate-10000"}},
		//},
		//{
		//	"list all",
		//	args{[]string{os.Args[0], "list", "--namespace", ""}},
		//},
		//{
		//	"list fate-10000",
		//	args{[]string{os.Args[0], "list", "--namespace", "fate-10000"}},
		//},
		//{
		//	"get",
		//	args{[]string{os.Args[0], "get", "--name", "fate-10000", "--namespace", "fate-10000"}},
		//},
		//{
		//	"delete",
		//	args{[]string{os.Args[0], "delete", "--name", "fate-10000", "--namespace", "fate-10000"}},
		//},
		//{
		//	"list",
		//	args{[]string{os.Args[0], "list", "--namespace", ""}},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.Args)
		})
	}
}
