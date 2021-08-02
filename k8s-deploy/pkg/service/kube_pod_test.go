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

package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service/kube"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type podtest struct {
	kube.Kube
}

var testpod = &v1.PodList{Items: []v1.Pod{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod-name",
			Namespace: "test-pod-namespace",
			Labels:    map[string]string{"name": "test-name"},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{Name: "test-pod-container"}, {Name: "test-pod-container-0"}},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod-name-1",
			Namespace: "test-pod-namespace-1",
			Labels:    map[string]string{"name": "test-name"},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{Name: "test-pod-container-2"}, {Name: "test-pod-container-3"}},
		},
	},
}}

func (e *podtest) GetPods(namespace, LabelSelector string) (*v1.PodList, error) {
	return testpod, nil
}

type podtesterr struct {
	kube.Kube
}

func (e *podtesterr) GetPods(namespace, LabelSelector string) (*v1.PodList, error) {
	return nil, errors.New("")
}

type podtestNoFind struct {
	kube.Kube
}

func (e *podtestNoFind) GetPods(namespace, LabelSelector string) (*v1.PodList, error) {
	return &v1.PodList{}, nil
}
func TestGetPodNameByModule(t *testing.T) {

	type args struct {
		namespace string
		name      string
		modules   string
		client    kubeClient
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "error",
			args: args{
				client: &podtesterr{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no find",
			args: args{
				client: &podtestNoFind{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				name:      testpod.Items[0].Labels["name"],
				namespace: testpod.Items[0].Namespace,
				modules:   testpod.Items[0].Spec.Containers[1].Name,
				client:    &podtest{},
			},
			want:    testpod.Items[0].Name,
			wantErr: false,
		},
		{
			name: "Success-1",
			args: args{
				name:      testpod.Items[0].Labels["name"],
				namespace: testpod.Items[0].Namespace,
				modules:   testpod.Items[1].Spec.Containers[0].Name,
				client:    &podtest{},
			},
			want:    testpod.Items[1].Name,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			KubeClient = tt.args.client
			got, err := GetPodNameByModule(tt.args.namespace, tt.args.name, tt.args.modules)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPodNameByModule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPodNameByModule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPodList(t *testing.T) {
	type args struct {
		name      string
		namespace string
		client    kubeClient
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Error",
			args: args{
				client: &podtesterr{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				name:      testpod.Items[0].Labels["name"],
				namespace: testpod.Items[0].Namespace,
				client:    &podtest{},
			},
			want:    []string{testpod.Items[0].Name, testpod.Items[1].Name},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			KubeClient = tt.args.client
			got, err := GetPodList(tt.args.name, tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPodList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPodList() = %v, want %v", got, tt.want)
			}
		})
	}
}
