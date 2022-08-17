package cli

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
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

	s5 = `
backend: eggroll
modules:
  - rollsite
  - clustermanager
python: 
`

	s6 = `
backend: spark_rabbitmq
modules:
  - pulsar
  - rabbitmq
rollsite: 
`
	s7 = `
modules:
  - mysql
  - python
  - fateboard
  - client
  - pulsar
  - hdfs

computing: Spark
federation: Eggroll
storage: HDFS
`

	s8 = `
modules:
  - pulsar
  - rabbitmq

computing: Eggroll
federation: Pulsar
storage: HDFS
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
		{"validateEmptyYaml", args{"", "", nil}, []error{SkipError("template or test yaml is empty")}},
		{"validateSameYaml", args{s1, s1, nil}, []error{ConfigError("computing error, not found"), ConfigError("the modules in your yaml is not valid")}},
		{"validateValidYaml", args{s1, s2, nil}, []error{ConfigError("computing error, not found"), ConfigError("the modules in your yaml is not valid")}},
		{"validateNotValidYaml", args{s1, s3, nil}, []error{ConfigError("computing error, not found"), ConfigError("the modules in your yaml is not valid"), ConfigError("your yaml at '/a/d', line 8 \n  'd: 3' may be redundant")}},
		{"validateYamlWithskippedKeys", args{s1, s4, []string{"d"}}, []error{ConfigError("computing error, not found"), ConfigError("the modules in your yaml is not valid")}},
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

func Test_getModules(t *testing.T) {
	type args struct {
		yamlMap map[string]interface{}
	}
	m4, _ := bufferToMap([]byte(s4))
	m5, _ := bufferToMap([]byte(s5))
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"no modules", args{m4}, nil, true},
		{"with modules", args{m5}, []string{"rollsite", "clustermanager"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getModules(tt.args.yamlMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("getModules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getModules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBackend(t *testing.T) {
	type args struct {
		yamlMap map[string]interface{}
	}

	m6, _ := bufferToMap([]byte(s6))
	m7, _ := bufferToMap([]byte(s7))
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{"no backends", args{m6}, nil, true},
		{"With backends", args{m7}, map[string]string{
			"computing":  "Spark",
			"federation": "Eggroll",
			"storage":    "HDFS",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBackend(tt.args.yamlMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBackend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getBackend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkCommonModules(t *testing.T) {
	type args struct {
		modules []string
	}
	common := []string{"mysql", "python", "fateboard", "client"}
	errs := []error{}
	for _, c := range common {
		errs = append(errs, ConfigError(fmt.Sprintf("common module %s not enabled", c)))
	}
	tests := []struct {
		name     string
		args     args
		wantErrs []error
	}{
		{"no backends", args{nil}, errs},
		{"full backends", args{common}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrs := checkCommonModules(tt.args.modules); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("checkCommonModules() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func Test_checkModuleBackend(t *testing.T) {
	type args struct {
		modules []string
		backend map[string]string
	}
	m7, _ := bufferToMap([]byte(s7))
	m8, _ := bufferToMap([]byte(s8))
	backend7, _ := getBackend(m7)
	backend8, _ := getBackend(m8)
	module7, _ := getModules(m7)
	module8, _ := getModules(m8)
	tests := []struct {
		name     string
		args     args
		wantErrs []error
	}{
		{"7", args{module7, backend7}, []error{ConfigError("module pulsar shall work with federation Pulsar but Eggroll")}},
		{"8", args{module8, backend8}, []error{ConfigError("module rabbitmq shall work with federation RabbitMQ but Pulsar")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrs := checkModuleBackend(tt.args.modules, tt.args.backend); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("checkModuleBackend() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func Test_checkComputing(t *testing.T) {
	type args struct {
		backend map[string]string
		modules []string
	}
	m7, _ := bufferToMap([]byte(s7))
	m8, _ := bufferToMap([]byte(s8))
	backend7, _ := getBackend(m7)
	backend8, _ := getBackend(m8)
	module7, _ := getModules(m7)
	module8, _ := getModules(m8)
	tests := []struct {
		name     string
		args     args
		wantErrs []error
	}{
		{"7", args{backend7, module7}, []error{ConfigError("computing Spark shall work with module spark"), ConfigError("computing Spark shall work with module nginx")}},
		{"8", args{backend8, module8}, []error{ConfigError("computing Eggroll shall work with module rollsite"), ConfigError("computing Eggroll shall work with module clustermanager"), ConfigError("computing Eggroll shall work with module nodemanager")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrs := checkComputing(tt.args.backend, tt.args.modules); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("checkComputing() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func Test_checkFederation(t *testing.T) {
	type args struct {
		backend map[string]string
		modules []string
	}
	m7, _ := bufferToMap([]byte(s7))
	m8, _ := bufferToMap([]byte(s8))
	backend7, _ := getBackend(m7)
	backend8, _ := getBackend(m8)
	module7, _ := getModules(m7)
	module8, _ := getModules(m8)
	tests := []struct {
		name     string
		args     args
		wantErrs []error
	}{
		{"7", args{backend7, module7}, []error{ConfigError("federation Eggroll shall work with module rollsite"), ConfigError("federation Eggroll shall work with module clustermanager"), ConfigError("federation Eggroll shall work with module nodemanager")}},
		{"8", args{backend8, module8}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrs := checkFederation(tt.args.backend, tt.args.modules); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("checkFederation() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func Test_checkStorage(t *testing.T) {
	type args struct {
		backend map[string]string
		modules []string
	}
	m7, _ := bufferToMap([]byte(s7))
	m8, _ := bufferToMap([]byte(s8))
	backend7, _ := getBackend(m7)
	backend8, _ := getBackend(m8)
	module7, _ := getModules(m7)
	module8, _ := getModules(m8)
	tests := []struct {
		name     string
		args     args
		wantErrs []error
	}{
		{"7", args{backend7, module7}, nil},
		{"8", args{backend8, module8}, []error{ConfigError("storage HDFS shall work with module hdfs")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrs := checkStorage(tt.args.backend, tt.args.modules); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("checkStorage() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func TestContainsSkipError(t *testing.T) {
	type args struct {
		errs []error
	}
	e1 := errors.New("")
	e2 := ConfigError("")
	e3 := SkipError("")
	e4 := fmt.Errorf("error :%w ", e3)
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"without skipErrors", args{[]error{e1, e2}}, false},
		{"with skipErrors", args{[]error{e3}}, true},
		{"with wrappedSkipErrors", args{[]error{e4}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsSkipError(tt.args.errs); got != tt.want {
				t.Errorf("ContainsSkipError() = %v, want %v", got, tt.want)
			}
		})
	}
}
