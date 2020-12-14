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

package kube

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

func TestGetPodLog(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the fake client.
	client := fake.NewSimpleClientset()

	// We will create an informer that writes added pods to a channel.
	pods := make(chan *v1.Pod, 1)
	informers := informers.NewSharedInformerFactory(client, 0)
	podInformer := informers.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			t.Logf("pod added: %s/%s", pod.Namespace, pod.Name)

			pod.Status.Conditions[0].Message = "aaab"

			pods <- pod
		},
	})

	podLister := informers.Core().V1().Pods().Lister()
	podLister.Pods("test-ns").Get("my-pod")
	// Make sure informers are running.
	informers.Start(ctx.Done())

	// This is not required in tests, but it serves as a proof-of-concept by
	// ensuring that the informer goroutine have warmed up and called List before
	// we send any events to it.
	cache.WaitForCacheSync(ctx.Done(), podInformer.HasSynced)

	// Inject an event into the fake client.
	p := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "my-pod"}}
	_, err := client.CoreV1().Pods("test-ns").Create(context.TODO(), p, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error injecting pod add: %v", err)
	}

	pod, err := client.CoreV1().Pods("test-ns").Get(context.TODO(), p.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("error injecting pod add: %v", err)
	}

	fmt.Println("aaa", pod.Name)
	fmt.Println("bbb", pod.Namespace)

	fmt.Println("bbb", pod.Status.String())

	res := client.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{})
	fmt.Println(res)
	logReader, err := res.Stream(context.TODO())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintf(os.Stdout, "POD LOGS: %s\n", pod.Name)
	_, err = io.Copy(os.Stdout, logReader)
	fmt.Fprintln(os.Stdout)
	select {
	case pod := <-pods:
		t.Logf("Got pod from channel: %s/%s", pod.Namespace, pod.Name)
	case <-time.After(wait.ForeverTestTimeout):
		t.Error("Informer did not get the added pod")
	}

}

func TestKube_GetPodLogs(t *testing.T) {
	// GetPodLogs GetPodLogs
	rest := KUBE.client.CoreV1().Pods("fate-9999").GetLogs("python-66fd9c4cb5-lw2sc",
		&corev1.PodLogOptions{Container: "python", Follow: true})

	podLogs, err := rest.Stream(KUBE.ctx)
	if err != nil {
		panic(err)
	}
	defer podLogs.Close()

	r := bufio.NewReader(podLogs)

	for {
		bytes, err := r.ReadByte()

		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			fmt.Println("io.EOF")
			return
		}

		fmt.Printf("%s", string(bytes))

	}

}
