package cli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"sigs.k8s.io/yaml"
)

type any interface{}
type KeyValue map[string]*TreeNode
type ListValue []*TreeNode

type TreeNode struct {
	leaf     bool
	lineno   int
	route    []string
	children KeyValue
	value    any
}

type ValidationTree struct {
	root                *TreeNode
	yamlMap, rawYamlMap map[string]interface{}
	lines               []string
}

type ValidationManager struct {
	templateTree, testTree *ValidationTree
	skippedKeys            []string
	preprocessor           []func(m *ValidationManager) []error
}

// SkipError is a type of error shows
// that you shall skip the validation
type SkipError string

// ConfigError is config validation error
type ConfigError string

func (e SkipError) Error() string   { return string(e) }
func (e ConfigError) Error() string { return string(e) }

// // trimComments trims the comments started with "# ".
func trimComments(t []byte) []byte {
	pattern := regexp.MustCompile(`^ *# `)
	for {
		ok := pattern.Match(t)
		if ok {
			t = bytes.Replace(t, []byte("# "), nil, 1)
		} else {
			break
		}
	}
	return t
}

// deconstructKey deconstructs the key to the original key and the lineno.
func deconstructKey(k interface{}) (string, int) {
	key, lineno := "", 0
	var err error
	switch k := k.(type) {
	case string:
		group := strings.Split(k, "__lineno__")
		if len(group) == 1 {
			key, lineno = k, 0
		} else {
			key = group[0]
			lineno, err = strconv.Atoi(group[1])
			if err != nil {
				lineno = 0
			}
		}
	default:
		key, lineno = fmt.Sprint(k), 0
	}
	return key, lineno
}

// NewTreeNode return default TreeNode.
func NewTreeNode() *TreeNode {
	node := new(TreeNode)
	node.leaf = false
	node.children = make(map[string]*TreeNode)
	return node
}

// mapToTreeNode recursively converts the yaml map to TreeNode,
// the route is the path to the current node.
// if node is a anomymous member in one array, the current route is @ArrayItem.
// node.value depends on the type of the key (map, list or basic type).
func mapToTreeNode(body interface{}, route []string) *TreeNode {
	node := NewTreeNode()
	node.route = route
	if body == nil {
		node.leaf = true
		node.value = nil
		return node
	}
	bodyType, bodyValue := reflect.TypeOf(body), reflect.ValueOf(body)
	switch bodyType.Kind() {
	case reflect.Map:
		for _, k := range bodyValue.MapKeys() {
			v := bodyValue.MapIndex(k).Interface()
			key, lineno := deconstructKey(k.Interface())
			r := append(route, key)
			child := mapToTreeNode(v, r)
			child.lineno = lineno
			node.children[key] = child
		}
		node.value = node.children
	case reflect.Array, reflect.Slice:
		arr := ListValue{}
		for i := 0; i < bodyValue.Len(); i++ {
			v := bodyValue.Index(i).Interface()
			arr = append(arr, mapToTreeNode(v, append(route, "@ArrayItem")))
		}
		node.value = arr
	default:
		node.leaf = true
		node.value = body
	}
	return node
}

// Contains checks whether an element is in slice/array/map.
func Contains(element interface{}, set interface{}) bool {
	setVal := reflect.ValueOf(set)
	switch setVal.Type().Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < setVal.Len(); i++ {
			if setVal.Index(i).Interface() == element {
				return true
			}
		}
	case reflect.Map:
		if setVal.MapIndex(reflect.ValueOf(element)).IsValid() {
			return true
		}
	default:
		panic(fmt.Sprintf("invalid type %v ", setVal.Type().Kind()))
	}

	return false
}

// GetValueTemplateExample gets the value template example from api.
func GetValueTemplateExample(chartName, chartVersion string) (value string, err error) {
	if !versionValid(chartVersion, []int{1, 9, 0}) {
		err = SkipError(fmt.Sprintf("version of %s does not meet the validation"+
			" requirement that chartVersion >= %s", chartVersion, "1.9.0"))
		return
	}

	c := new(Chart)
	req := &Request{
		Type: "GET",
		Path: fmt.Sprintf("chart/valueTemplateExample/%s/%s", chartName, chartVersion),
		Body: nil,
	}
	msg, err := GetResponse(c, req, VALUE)
	if err != nil {
		err = SkipError(fmt.Sprintf("get value template example for validation failed %v,"+
			", you may need to upgrade KubeFATE service", err))
		return
	}

	if result, ok := msg.(*ValueResult); ok {
		value = result.Data
	}

	if result, ok := msg.(*ChartResultErr); ok {
		err = SkipError(fmt.Sprintf("get value template example failed\n %v", result.Error))
	}
	return
}

// compareTwoTrees recursively compares the two trees, walking through the nodes not skipped.
func compareTwoTrees(rootTemp, rootTest *TreeNode,
	testLines, skippedKeys []string) (errs []error) {
	valueTemp, valueTest := rootTemp.value, rootTest.value
	typeTemp, typeTest := reflect.TypeOf(valueTemp), reflect.TypeOf(valueTest)
	if typeTemp != typeTest {
		route := strings.Join(rootTest.route, "/")
		errs = append(errs,
			ConfigError(fmt.Sprintf("your yaml at '%s', line %d \n  '%s' may not match the type",
				route, rootTest.lineno, testLines[rootTest.lineno])))
		return
	}
	switch valueTest := valueTest.(type) {
	case KeyValue:
		for k, v := range valueTest {
			if Contains(k, skippedKeys) {
				continue
			}
			if childTemp, ok := valueTemp.(KeyValue)[k]; !ok {
				route := strings.Join(v.route, "/")
				errs = append(errs,
					ConfigError(fmt.Sprintf("your yaml at '%s', line %d \n  '%s' may be redundant",
						route, v.lineno, testLines[v.lineno])))
			} else {
				errs = append(errs, compareTwoTrees(childTemp, v, testLines, skippedKeys)...)
			}
		}
	case ListValue:
		item := valueTemp.(ListValue)[0]
		for _, v := range valueTest {
			compareTwoTrees(item, v, testLines, skippedKeys)
		}
	default:
	}
	return
}

func (m *ValidationManager) Validate() error {
	compareTwoTrees(m.templateTree.root, m.testTree.root, m.testTree.lines, m.skippedKeys)
	return nil
}

// versionValid checks if the chart version is valid.
func versionValid(chartVersion string, startVersion []int) (valid bool) {
	chartVersion = strings.TrimLeft(chartVersion, "v")
	for i, v := range strings.Split(chartVersion, ".") {
		if i >= len(startVersion) {
			return
		}
		if v, err := strconv.Atoi(v); err != nil || v < startVersion[i] {
			return
		} else if v > startVersion[i] {
			break
		}
	}
	valid = true
	return
}

func (m *ValidationManager) compareTwoTrees() []error {
	return compareTwoTrees(m.templateTree.root, m.testTree.root, m.testTree.lines, m.skippedKeys)
}

// yamlStringToBuffer reads the yaml value and returns
// the content (may be modified) in []byte,the original lines in []string and error.
func yamlStringToBuffer(value string, restoreComments, markLineno bool) ([]byte, []string, error) {
	reader := bufio.NewReader(strings.NewReader(value))
	buffer := make([]byte, 0, 10)
	// 1 for the first lineno
	lines := make([]string, 1, 10)

	linenoReg := regexp.MustCompile(`:`)
	for lineno := 1; ; lineno++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, nil, err
			}
		}
		if !Contains(byte(':'), line) && !Contains(byte('-'), line) {
			// if one line is pure comment, treat it as a blank line.
			line = []byte("")
		}

		lines = append(lines, strings.TrimSpace(string(line)))
		if restoreComments {
			line = trimComments(line)
		}
		if markLineno {
			line = linenoReg.ReplaceAll(line, []byte(fmt.Sprintf("__lineno__%d:", lineno)))
		}
		buffer = append(buffer, append(line, '\n')...)
	}
	return buffer, lines, nil
}

// bufferToMap unmarshals the byte buffer to map,
// keys are string and numbers are float64.
func bufferToMap(buffer []byte) (m map[string]interface{}, err error) {
	err = yaml.Unmarshal(buffer, &m)
	return
}

// buildValidationTree builds the validation tree from yaml string.
func buildValidationTree(yamlString string, restoreComments,
	markLineno bool) (*ValidationTree, error) {
	yamlBuffer, lines, err := yamlStringToBuffer(yamlString, restoreComments, markLineno)
	if err != nil {
		return nil, err
	}
	yamlMap, err := bufferToMap(yamlBuffer)
	if err != nil {
		return nil, err
	}
	rawYamlMap, err := bufferToMap([]byte(yamlString))
	if err != nil {
		return nil, err
	}
	root := mapToTreeNode(yamlMap, []string{""})
	return &ValidationTree{root, yamlMap, rawYamlMap, lines}, nil
}

// getSkippedKeys returns the skippedKeys in a map.
// The type of skippedKeys array is []interface{},
// so we need to transform every element to string.
func getSkippedKeys(m map[string]interface{}) (skippedKeys []string) {
	value, ok := m["skippedKeys"]
	if !ok {
		return
	}
	slice, ok := value.([]interface{})
	if !ok {
		return
	}
	for _, v := range slice {
		skippedKeys = append(skippedKeys, v.(string))
	}
	return skippedKeys
}

// getModules extracts modules []string from map
func getModules(yamlMap map[string]interface{}) ([]string, error) {
	modules, ok := yamlMap["modules"].([]interface{})
	if !ok {
		return nil, ConfigError("the modules in your yaml is not valid")
	}
	var s []string
	for _, module := range modules {
		m := module.(string)
		s = append(s, m)
	}
	return s, nil
}

// getModules extracts backend components from map
func getBackend(yamlMap map[string]interface{}) (map[string]string, error) {
	m := map[string]string{}
	if computing, ok := yamlMap["computing"].(string); ok {
		m["computing"] = computing
	} else {
		return nil, ConfigError("computing error, not found")
	}
	if federation, ok := yamlMap["federation"].(string); ok {
		m["federation"] = federation
	} else {
		return nil, ConfigError("federation error, not found")
	}
	if storage, ok := yamlMap["storage"].(string); ok {
		m["storage"] = storage
	} else {
		return nil, ConfigError("storage error, not found")
	}
	return m, nil
}

// checkCommonModules checks common modules' presence
func checkCommonModules(modules []string) (errs []error) {
	common := []string{"mysql", "python", "fateboard", "client"}
	for _, c := range common {
		if !Contains(c, modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("common module %s not enabled", c)))
		}
	}
	return
}

// checkModuleBackend checks each module has correct backend mapping
func checkModuleBackend(modules []string, backend map[string]string) (errs []error) {
	for _, m := range modules {
		switch m {
		case "spark":
			if c := backend["computing"]; c != "Spark" {
				errs = append(errs,
					ConfigError("module spark shall work with computing Spark but "+c))
			}
		case "hdfs":
			if s := backend["storage"]; s != "HDFS" {
				errs = append(errs,
					ConfigError("module hdfs shall work with computing HDFS but "+s))
			}
		case "nginx":
			if c := backend["computing"]; !Contains(c, []string{"Spark", "Spark_local"}) {
				errs = append(errs,
					ConfigError("module nginx shall work with computing Spark/Spark_local but "+c))
			}
		case "pulsar":
			if f := backend["federation"]; f != "Pulsar" {
				errs = append(errs,
					ConfigError("module pulsar shall work with federation Pulsar but "+f))
			}
		case "rabbitmq":
			if f := backend["federation"]; f != "RabbitMQ" {
				errs = append(errs,
					ConfigError("module rabbitmq shall work with federation RabbitMQ but "+f))
			}
		}
	}
	return
}

// moduleValidator validates yaml from modules' view
func moduleValidator(m *ValidationManager) (errs []error) {
	yamlMap := m.testTree.rawYamlMap
	if chartName, err := getChartName(yamlMap); err != nil || chartName != "fate" {
		// if chart is not fate, just skip this validation
		return []error{}
	}
	backend, err := getBackend(yamlMap)
	if err != nil {
		return []error{err}
	}
	modules, err := getModules(yamlMap)
	if err != nil {
		return []error{err}
	}
	errs = append(errs, checkCommonModules(modules)...)
	errs = append(errs, checkModuleBackend(modules, backend)...)
	return
}

// checkComputing checks computing has correct modules enabled
func checkComputing(backend map[string]string, modules []string) (errs []error) {
	key := "computing"
	switch c := backend[key]; c {
	case "Eggroll":
		if !Contains("rollsite", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, c, "rollsite")))
		}
		if !Contains("clustermanager", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, c, "clustermanager")))
		}
		if !Contains("nodemanager", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, c, "nodemanager")))
		}
	case "Spark":
		if !Contains("spark", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, c, "spark")))
		}
		if !Contains("nginx", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, c, "nginx")))
		}
	case "Spark_local":
		if !Contains("nginx", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, c, "nginx")))
		}
	default:
		errs = append(errs, ConfigError(fmt.Sprintf("%s module %s is not supported", key, c)))
	}
	return
}

// checkFederation checks federation has correct modules enabled
func checkFederation(backend map[string]string, modules []string) (errs []error) {
	key := "federation"
	switch f := backend[key]; f {
	case "Eggroll":
		if !Contains("rollsite", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, f, "rollsite")))
		}
		if !Contains("clustermanager", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, f, "clustermanager")))
		}
		if !Contains("nodemanager", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, f, "nodemanager")))
		}
	case "Pulsar":
		if !Contains("pulsar", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, f, "pulsar")))
		}
	case "RabbitMQ":
		if !Contains("rabbitmq", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, f, "rabbitmq")))
		}
	default:
		errs = append(errs, ConfigError(fmt.Sprintf("%s module %s is not supported", key, f)))
	}
	return
}

// checkStorage checks storage has correct modules enabled
func checkStorage(backend map[string]string, modules []string) (errs []error) {
	key := "storage"
	switch s := backend[key]; s {
	case "Eggroll":
		if !Contains("rollsite", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, s, "rollsite")))
		}
		if !Contains("clustermanager", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, s, "clustermanager")))
		}
		if !Contains("nodemanager", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, s, "nodemanager")))
		}
	case "HDFS":
		if !Contains("hdfs", modules) {
			errs = append(errs, ConfigError(fmt.Sprintf("%s %s shall work with module %s",
				key, s, "hdfs")))
		}
	case "LocalFS":
	default:
		errs = append(errs, ConfigError(fmt.Sprintf("%s module %s is not supported", key, s)))
	}
	return
}

func getChartName(m map[string]interface{}) (string, error) {
	chartName, ok := m["chartName"]
	err := errors.New("chart name not found")
	if !ok {
		return "", err
	}
	name, ok := chartName.(string)
	if !ok || name != "fate" {
		return "", err
	}
	return name, nil
}

// backendValidator validates yaml from backend' view
func backendValidator(m *ValidationManager) (errs []error) {
	yamlMap := m.testTree.rawYamlMap
	if chartName, err := getChartName(yamlMap); err != nil || chartName != "fate" {
		// if chart is not fate, just skip this validation
		return []error{}
	}
	modules, err := getModules(yamlMap)
	if err != nil {
		return []error{err}
	}

	backend, err := getBackend(yamlMap)
	if err != nil {
		return []error{err}
	}
	errs = append(errs, checkComputing(backend, modules)...)
	errs = append(errs, checkFederation(backend, modules)...)
	errs = append(errs, checkStorage(backend, modules)...)
	return
}

// ContainsSkipError returns if error list has SkipError
func ContainsSkipError(errs []error) bool {
	var skip SkipError
	for _, e := range errs {
		if errors.As(e, &skip) {
			return true
		}
	}
	return false
}

// RegisterPreprocessor registers a callback
func (m *ValidationManager) RegisterPreprocessor(p func(m *ValidationManager) []error) {
	m.preprocessor = append(m.preprocessor, p)
}

// preprocess runs callbacks
func (m *ValidationManager) preprocess() (errs []error) {
	for _, p := range m.preprocessor {
		errs = append(errs, p(m)...)
	}
	return
}

// ValidateYaml validates the yaml file.
func ValidateYaml(templateValue, testValue string, skippedKeys []string,
	preprocessors ...func(m *ValidationManager) []error) (errs []error) {
	if templateValue == "" || testValue == "" {
		return []error{SkipError("template or test yaml is empty")}
	}
	templateTree, err := buildValidationTree(templateValue, true, false)
	if err != nil {
		return []error{err}
	}
	testTree, err := buildValidationTree(testValue, false, true)
	if err != nil {
		return []error{err}
	}
	m := &ValidationManager{templateTree, testTree, skippedKeys, preprocessors}
	if m.templateTree == nil || m.testTree == nil {
		return []error{ConfigError("building validation tree failed")}
	}
	m.RegisterPreprocessor(moduleValidator)
	m.RegisterPreprocessor(backendValidator)
	errs = append(errs, m.preprocess()...)
	errs = append(errs, m.compareTwoTrees()...)
	return
}
