package service

import (
	"fmt"
	"runtime"
)

var (
	kubefateVersion = "v0.0.0"
	gitVersion      = "v0.0.0-master+$Format:%h$"
	gitCommit       = "$Format:%H$"          // sha1 from git, output of $(git rev-parse HEAD)
	buildDate       = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

func GetVersion() *Version {
	return &Version{
		KubefateVersion: kubefateVersion,
		GitVersion:      gitVersion,
		GitCommit:       gitCommit,
		BuildDate:       buildDate,
		GoVersion:       runtime.Version(),
		Compiler:        runtime.Compiler,
		Platform:        fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

type Version struct {
	KubefateVersion string `json:"kubefateVersion"`
	GitVersion      string `json:"gitVersion"`
	GitCommit       string `json:"gitCommit"`
	BuildDate       string `json:"buildDate"`
	GoVersion       string `json:"goVersion"`
	Compiler        string `json:"compiler"`
	Platform        string `json:"platform"`
}
