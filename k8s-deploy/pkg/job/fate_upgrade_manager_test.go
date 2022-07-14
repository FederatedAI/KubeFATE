package job

import "testing"

func Test_ConstructFumClusterData(t *testing.T) {
	actual := ConstructFumClusterData("fate", "fate_dev", "1.8.0", "1.9.0")
	println(string(actual))
}
