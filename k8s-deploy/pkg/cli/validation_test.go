package cli

import (
	"errors"
	"reflect"
	"testing"
)

const (
	_templateFilename = "../../../helm-charts/FATE/values-template-example.yaml"
	_testFilename     = "../../examples/party-10000/cluster.yaml"
)

var (
	s0 = `
chartName: fate
chartVersion: 1.9.0
`

	s1 = `
chartName: fate
chartVersion: 1.9.0
skippedKeys:
  - str
a:
  b: 2
  c: [3, 4]
`

	s2 = `
chartName: fate
chartVersion: 1.9.0
skippedKeys:
  - str
a:
  b: 2
`

	s3 = `
chartName: fate
chartVersion: 1.9.0
skippedKeys:
  - str
a:
  b: 2
  d: 3
`

	s4 = `
chartName: fate
chartVersion: 1.9.0
skippedKeys:
  - d
a:
  b: 2
  d: 3
`
)

var (
	m0 = map[string]interface{}{
		"chartName":    "fate",
		"chartVersion": "1.9.0",
	}
	m1 = map[string]interface{}{
		"chartName":    "fate",
		"chartVersion": "1.9.0",
		"skippedKeys":  []interface{}{"str"},
		"a": map[string]interface{}{
			"b": float64(2),
			"c": []interface{}{float64(3), float64(4)},
		},
	}
	m2 = map[string]interface{}{
		"chartName":    "fate",
		"chartVersion": "1.9.0",
		"skippedKeys":  []interface{}{"str"},
		"a": map[string]interface{}{
			"b": float64(2),
		},
	}
	m3 = map[string]interface{}{
		"chartName":    "fate",
		"chartVersion": "1.9.0",
		"skippedKeys":  []interface{}{"str"},
		"a": map[string]interface{}{
			"b": float64(2),
			"d": float64(3),
		},
	}
	m4 = map[string]interface{}{
		"chartName":    "fate",
		"chartVersion": "1.9.0",
		"skippedKeys":  []interface{}{"d"},
		"a": map[string]interface{}{
			"b": float64(2),
			"d": float64(3),
		},
	}
)

func Test_trimComments(t *testing.T) {
	type args struct {
		t []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"trimComment", args{[]byte("# comment\n")}, []byte("comment\n")},
		{"trimComments", args{[]byte("   # # # comment")}, []byte("   comment")},
		{"trimLeftComments", args{[]byte("   # # # comment # f")}, []byte("   comment # f")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimComments(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trimComments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deconstructKey(t *testing.T) {
	type args struct {
		k interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int
	}{
		{"deconstructKey", args{interface{}("key")}, "key", 0},
		{"deconstructKeyWithLineno", args{interface{}("key__lineno__21")}, "key", 21},
		{"deconstructNotValidKey", args{interface{}(21)}, "21", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := deconstructKey(tt.args.k)
			if got != tt.want {
				t.Errorf("deconstructKey() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("deconstructKey() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_versionValid(t *testing.T) {
	type args struct {
		chartVersion string
		startVersion []int
	}
	tests := []struct {
		name      string
		args      args
		wantValid bool
	}{
		{"versionEqual", args{"1.0.0", []int{1, 0, 0}}, true},
		{"versionLower", args{"1.0.0", []int{1, 1, 0}}, false},
		{"versionHigher", args{"1.2.0", []int{1, 1, 0}}, true},
		{"versionWithV", args{"v1.2.0", []int{1, 1, 0}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotValid := versionValid(tt.args.chartVersion, tt.args.startVersion); gotValid != tt.wantValid {
				t.Errorf("versionValid() = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}

func TestGetValueTemplateExample(t *testing.T) {
	type args struct {
		chartName    string
		chartVersion string
	}
	tests := []struct {
		name      string
		args      args
		wantValue string
		wantErr   bool
	}{
		{"getLowVersionFateValueTemplateExample", args{"fate", "1.8.0"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, err := GetValueTemplateExample(tt.args.chartName, tt.args.chartVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValueTemplateExample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotValue != tt.wantValue {
				t.Errorf("GetValueTemplateExample() = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func Test_bufferToMap(t *testing.T) {
	type args struct {
		buffer []byte
	}
	tests := []struct {
		name    string
		args    args
		wantM   map[string]interface{}
		wantErr bool
	}{
		{"bufferToMap0", args{[]byte(s0)}, m0, false},
		{"bufferToMap1", args{[]byte(s1)}, m1, false},
		{"bufferToMap2", args{[]byte(s2)}, m2, false},
		{"bufferToMap3", args{[]byte(s3)}, m3, false},
		{"bufferToMap4", args{[]byte(s4)}, m4, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, err := bufferToMap(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("bufferToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotM, tt.wantM) {
				t.Errorf("bufferToMap() = %v, want %v", gotM, tt.wantM)
			}
		})
	}
}

func TestValidateYaml(t *testing.T) {
	type args struct {
		templateValue string
		testValue     string
		skippedKeys   []string
	}
	tests := []struct {
		name     string
		args     args
		wantErrs []error
	}{
		{"validateEmptyYaml", args{"", "", nil}, []error{errors.New("template or test yaml is empty")}},
		{"validateSameYaml", args{s1, s1, nil}, []error{}},
		{"validateValidYaml", args{s1, s2, nil}, []error{}},
		{"validateNotValidYaml", args{s1, s3, nil}, []error{errors.New("Your yaml at '/a/d', line 8 \n  'd: 3' may be redundant\n")}},
		{"validateYamlWithskippedKeys", args{s1, s4, []string{"d"}}, []error{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrs := ValidateYaml(tt.args.templateValue, tt.args.testValue, tt.args.skippedKeys); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				if len(gotErrs) == 0 && len(tt.wantErrs) == 0 {
					return
				}
				t.Errorf("ValidateYaml() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func Test_getSkippedKeys(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	m0, _ := bufferToMap([]byte(s0))
	m4, _ := bufferToMap([]byte(s4))
	tests := []struct {
		name            string
		args            args
		wantSkippedKeys []string
	}{
		{"getSkippedKeys", args{m0}, nil},
		{"getEmptySkippedKeys", args{m4}, []string{"d"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSkippedKeys := getSkippedKeys(tt.args.m); !reflect.DeepEqual(gotSkippedKeys, tt.wantSkippedKeys) {
				t.Errorf("getSkippedKeys() = %v, want %v", gotSkippedKeys, tt.wantSkippedKeys)
			}
		})
	}
}
