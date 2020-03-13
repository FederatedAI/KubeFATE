package service

import (
	"fmt"
	"testing"

	"k8s.io/client-go/rest"
)

func TestAbc(t *testing.T) {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config.String())
}

func Test_checkClusterStatus(t *testing.T) {
	type args struct {
		namespace string
		name      string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "cluster is running",
			args: args{
				namespace: "fate-10000",
				name:      "fate-10000",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckClusterStatus(tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkClusterStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkClusterStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPods(t *testing.T) {
	type args struct {
		name string
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				name: "fate-10000",
				namespace: "fate-10000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPods(tt.args.name,tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("Namespace, Name, Status")
			for _, v := range got.Items {
				fmt.Printf("%s, %s, %s\n", v.Namespace, v.Name, v.Status.Phase)
			}
		})
	}
}
