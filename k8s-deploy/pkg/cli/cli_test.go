package cli

import (
	"fate-cloud-agent/pkg/utils/logging"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	type args struct {
		Args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"help",
			args{[]string{os.Args[0], "help"}},
		},
		{
			"cluster -help",
			args{[]string{os.Args[0], "cluster", "--help"}},
		},
		{
			"cluster list",
			args{[]string{os.Args[0], "cluster", "list"}},
		},
		//{
		//	"get -help",
		//	args{[]string{os.Args[0], "get", "--help"}},
		//},
		//{
		//	"delete -help",
		//	args{[]string{os.Args[0], "delete", "--help"}},
		//},
		//{
		//	"list all",
		//	args{[]string{os.Args[0], "list", "--namespace", ""}},
		//},
		//{
		//	"install",
		//	args{[]string{os.Args[0], "install", "--name", "fate-10000", "--namespace", "fate-10000", "--chart", "E:\\machenlong\\AI\\github\\owlet42\\KubeFATE\\k8s-deploy\\fate-10000"}},
		//},
		//{
		//	"list all",
		//	args{[]string{os.Args[0], "list", "--namespace", ""}},
		//},
		//{
		//	"list fate-10000",
		//	args{[]string{os.Args[0], "list", "--namespace", "fate-10000"}},
		//},
		//{
		//	"get",
		//	args{[]string{os.Args[0], "get", "--name", "fate-10000", "--namespace", "fate-10000"}},
		//},
		//{
		//	"delete",
		//	args{[]string{os.Args[0], "delete", "--name", "fate-10000", "--namespace", "fate-10000"}},
		//},
		//{
		//	"list",
		//	args{[]string{os.Args[0], "list", "--namespace", ""}},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.Args)
		})
	}
}
