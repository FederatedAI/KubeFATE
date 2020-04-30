package service

import (
	"reflect"
	"testing"
)

func TestGetIngressUrl(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				name:      "fate-serving-9999",
				namespace: "fate-serving-9999",
			},
			want:    []string{"9999.serving-proxy.kubefate.net"},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				name:      "fate-9999",
				namespace: "fate-9999",
			},
			want:    []string{"9999.fateboard.kubefate.net"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIngressUrl(tt.args.name, tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIngressUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIngressUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
