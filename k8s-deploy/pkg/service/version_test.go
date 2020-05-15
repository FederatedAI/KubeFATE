package service

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name string
		want *Version
	}{
		// TODO: Add test cases.
		{
			name: "GetVersion test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetVersion()
			t.Logf("GetVersion() = %v", got)
		})
	}
}
