package service

import (
	"reflect"
	"testing"
)

func TestGetProxySvcNodePorts(t *testing.T) {
	type args struct {
		namespace string
		name      string
	}
	tests := []struct {
		name    string
		args    args
		want    []int32
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				namespace: "fate-10000", name: "fate-10000",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProxySvcNodePorts(tt.args.name, tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProxySvc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProxySvc() = %v, want %v", got, tt.want)
			}
		})
	}
}
