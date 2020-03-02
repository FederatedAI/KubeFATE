package service

import (
	"reflect"
	"testing"
)

func TestGetNodeIp(t *testing.T) {
	tests := []struct {
		name    string
		want    []int32
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNodeIp()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNodeIp() = %v, want %v", got, tt.want)
			}
		})
	}
}
