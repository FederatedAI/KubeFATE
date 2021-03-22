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

package kube

import (
	"io"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Log interface
type Log interface {
	GetPodLogs(namespace, podName string, args *PodLogArgs) (io.ReadCloser, error)
}

type PodLogArgs struct {
	Container                    string
	Follow                       bool
	Previous                     bool
	SinceSeconds                 *int64
	SinceTime                    time.Time
	Timestamps                   bool
	TailLines                    *int64
	LimitBytes                   *int64
	InsecureSkipTLSVerifyBackend bool
}

// GetPodLogs GetPodLogs
func (e *Kube) GetPodLogs(namespace, podName string, args *PodLogArgs) (io.ReadCloser, error) {

	rest := e.client.CoreV1().Pods(namespace).GetLogs(podName, getPodLogOptions(args))

	podLogs, err := rest.Stream(e.ctx)
	if err != nil {
		return nil, err
	}
	//defer podLogs.Close()

	return podLogs, nil
}

func getPodLogOptions(podLogArgs *PodLogArgs) *corev1.PodLogOptions {
	return &corev1.PodLogOptions{
		Container:    podLogArgs.Container,
		Follow:       podLogArgs.Follow,
		Previous:     podLogArgs.Previous,
		SinceSeconds: podLogArgs.SinceSeconds,
		SinceTime: func() *metav1.Time {
			if podLogArgs.SinceTime.IsZero() || podLogArgs.SinceSeconds != nil {
				return nil
			}
			return &metav1.Time{
				Time: podLogArgs.SinceTime,
			}
		}(),
		Timestamps:                   podLogArgs.Timestamps,
		TailLines:                    podLogArgs.TailLines,
		LimitBytes:                   podLogArgs.LimitBytes,
		InsecureSkipTLSVerifyBackend: podLogArgs.InsecureSkipTLSVerifyBackend,
	}
}
