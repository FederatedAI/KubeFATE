package cli

import (
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"
	"os"
	"testing"
)

func TestCluster(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	type args struct {
		Args []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"cluster -help",
			args{[]string{os.Args[0], "cluster", "--help"}},
		},
		{
			"cluster list",
			args{[]string{os.Args[0], "cluster", "list"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.Args)
		})
	}
}
