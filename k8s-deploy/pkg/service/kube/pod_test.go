package kube

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

func TestKube_GetPods(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kube := KUBE
	kube.client = fake.NewSimpleClientset()
	fmt.Println("get client")
	informers := informers.NewSharedInformerFactory(kube.client, 0)
	podInformer := informers.Core().V1().Pods().Informer()

	informers.Start(ctx.Done())
	cache.WaitForCacheSync(ctx.Done(), podInformer.HasSynced)

	type args struct {
		name          string
		namespace     string
		LabelSelector string
	}
	tests := []struct {
		name    string
		e       *Kube
		args    args
		want    *corev1.PodList
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			e:    &kube,
			args: args{
				name:          "",
				namespace:     "",
				LabelSelector: "app=pod, lable=test",
			},
			want: &v1.PodList{
				Items: []v1.Pod{
					v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod-0", Namespace: "defaule", Labels: map[string]string{"app": "pod", "lable": "test", "name": "pod-0"}}},
					v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod-1", Namespace: "defaule", Labels: map[string]string{"app": "pod", "lable": "test", "name": "pod-1"}}},
					v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod-2", Namespace: "defaule", Labels: map[string]string{"app": "pod", "lable": "test", "name": "pod-2"}}},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// add a pods
			for _, v := range tt.want.Items {
				_, err := kube.client.CoreV1().Pods(v.Namespace).Create(context.TODO(), &v, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("error injecting pod add: %v", err)
				}
			}
			// add a pod of no lable
			_, err := kube.client.CoreV1().Pods("defaule").Create(context.TODO(), &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "nolable", Namespace: "defaule"}}, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("error injecting pod add: %v", err)
			}

			got, err := tt.e.GetPods(tt.args.namespace, tt.args.LabelSelector)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kube.GetPods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kube.GetPods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKube_GetPod(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kube := KUBE
	kube.client = fake.NewSimpleClientset()
	fmt.Println("get client")
	informers := informers.NewSharedInformerFactory(kube.client, 0)
	podInformer := informers.Core().V1().Pods().Informer()

	informers.Start(ctx.Done())
	cache.WaitForCacheSync(ctx.Done(), podInformer.HasSynced)

	type args struct {
		podName   string
		namespace string
	}
	tests := []struct {
		name    string
		e       *Kube
		args    args
		want    *corev1.Pod
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			e:    &kube,
			args: args{
				podName:   "python",
				namespace: "defaule",
			},
			want: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "python", Namespace: "defaule"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := kube.client.CoreV1().Pods(tt.args.namespace).Create(context.TODO(), tt.want, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("error injecting pod add: %v", err)
			}
			got, err := tt.e.GetPod(tt.args.podName, tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kube.GetPod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kube.GetPod() = %v, want %v", got, tt.want)
			}
		})
	}
}
