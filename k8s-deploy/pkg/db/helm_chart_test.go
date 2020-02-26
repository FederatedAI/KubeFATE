package db

import (
	"reflect"
	"testing"
)

func TestHelmChart_FindHelmByVersion(t *testing.T) {
	InitConfigForTest()
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    *HelmChart
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "read",
			args:    args{
				version: "v1.2.0",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindHelmByVersion(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("HelmChart.FindHelmByVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HelmChart.FindHelmByVersion() = %+v, want %v", got, tt.want)
			}
		})
	}
}
