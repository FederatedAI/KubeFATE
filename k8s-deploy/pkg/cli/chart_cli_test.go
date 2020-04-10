package cli

import (
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"
	"os"
	"testing"
)

func TestChartCreateCommand(t *testing.T) {

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
			"chart upload",
			args{[]string{os.Args[0], "chart", "upload", "-f", "X:/AI/owlet42/KubeFATE/k8s-deploy/docs/fate-party-1.3.0.tgz"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.Args)
		})
	}
}
