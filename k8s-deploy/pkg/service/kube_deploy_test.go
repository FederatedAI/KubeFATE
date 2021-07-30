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
	"testing"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestCheckDeploy(t *testing.T) {
	type args struct {
		deploy *v1.Deployment
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "True",
			args: args{
				deploy: &v1.Deployment{
					Status: v1.DeploymentStatus{
						Conditions: []v1.DeploymentCondition{
							{
								Type:   v1.DeploymentAvailable,
								Status: corev1.ConditionTrue,
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "nil",
			args: args{
				deploy: &v1.Deployment{},
			},
			want: false,
		},
		{
			name: "Zero",
			args: args{
				deploy: &v1.Deployment{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckDeploy(tt.args.deploy); got != tt.want {
				t.Errorf("CheckDeploy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDeploys(t *testing.T) {
	type args struct {
		deploys *v1.DeploymentList
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil",
			args: args{},
			want: false,
		},
		{
			name: "count-0",
			args: args{
				deploys: &v1.DeploymentList{},
			},
			want: false,
		},
		{
			name: "one-false",
			args: args{
				deploys: &v1.DeploymentList{
					Items: []v1.Deployment{
						{
							Status: v1.DeploymentStatus{
								Conditions: []v1.DeploymentCondition{
									{
										Type:   v1.DeploymentAvailable,
										Status: corev1.ConditionTrue,
									},
								},
							},
						},
						{
							Status: v1.DeploymentStatus{
								Conditions: []v1.DeploymentCondition{
									{
										Type:   v1.DeploymentAvailable,
										Status: corev1.ConditionFalse,
									},
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "one-not-Available",
			args: args{
				deploys: &v1.DeploymentList{
					Items: []v1.Deployment{
						{
							Status: v1.DeploymentStatus{
								Conditions: []v1.DeploymentCondition{
									{
										Type:   v1.DeploymentAvailable,
										Status: corev1.ConditionTrue,
									},
								},
							},
						},
						{
							Status: v1.DeploymentStatus{
								Conditions: []v1.DeploymentCondition{
									{
										Type:   v1.DeploymentProgressing,
										Status: corev1.ConditionTrue,
									},
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "all-Available",
			args: args{
				deploys: &v1.DeploymentList{
					Items: []v1.Deployment{
						{
							Status: v1.DeploymentStatus{
								Conditions: []v1.DeploymentCondition{
									{
										Type:   v1.DeploymentAvailable,
										Status: corev1.ConditionTrue,
									},
								},
							},
						},
						{
							Status: v1.DeploymentStatus{
								Conditions: []v1.DeploymentCondition{
									{
										Type:   v1.DeploymentAvailable,
										Status: corev1.ConditionTrue,
									},
								},
							},
						},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckDeploys(tt.args.deploys); got != tt.want {
				t.Errorf("CheckDeploys() = %v, want %v", got, tt.want)
			}
		})
	}
}
