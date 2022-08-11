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

	"github.com/fatih/color"
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
	root    *TreeNode
	yamlMap map[string]interface{}
	lines   []string
}

type ValidationManager struct {
	templateTree, testTree *ValidationTree
	skippedKeys            []string
}

type VersionNotValidError struct {
	Version    string
	LowerBound string
}

func (e VersionNotValidError) Error() string {
	return fmt.Sprintf("version of %s does not meet the validation requirement that chartVersion >= %s", e.Version, e.LowerBound)
}

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

// compareTwoTrees recursively compares the two trees, walking through the nodes not skipped.
func compareTwoTrees(rootTemp, rootTest *TreeNode,
	testLines, skippedKeys []string) (errs []error) {
	valueTemp, valueTest := rootTemp.value, rootTest.value
	typeTemp, typeTest := reflect.TypeOf(valueTemp), reflect.TypeOf(valueTest)
	if typeTemp != typeTest {
		route := strings.Join(rootTest.route, "/")
		errs = append(errs,
			fmt.Errorf("your yaml at '%s', line %d \n  '%s' may not match the type",
				route, rootTest.lineno, testLines[rootTest.lineno]))
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
					fmt.Errorf("your yaml at '%s', line %d \n  '%s' may be redundant",
						route, v.lineno, testLines[v.lineno]))
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

// GetValueTemplateExample gets the value template example from api.
func GetValueTemplateExample(chartName, chartVersion string) (value string, err error) {
	if !versionValid(chartVersion, []int{1, 9, 0}) {
		err = VersionNotValidError{chartVersion, "1.9.0"}
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
		err = fmt.Errorf("get value template example for validation failed\n %v, you may need to upgrade KubeFATE service", err)
		return
	}

	if result, ok := msg.(*ValueResult); ok {
		value = result.Data
	}

	if result, ok := msg.(*ChartResultErr); ok {
		err = fmt.Errorf("get value template example failed\n %v", errors.New(result.Error))
	}
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
	root := mapToTreeNode(yamlMap, []string{""})
	return &ValidationTree{root, yamlMap, lines}, nil
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

func alertUserIfModulesNotMatchBackend(yamlMap map[string]interface{}) {
	var allModules = []string{"rollsite", "clustermanager", "nodemanager", "mysql", "python",
		"fateboard", "client", "lbrollsite", "spark", "hdfs", "nginx", "rabbitmq", "pulsar"}
	var backendModules = map[string][]string{
		"eggroll": {"rollsite", "clustermanager", "nodemanager",
			"mysql", "python", "fateboard", "client"},
		"spark_rabbitmq": {"mysql", "python", "fateboard", "client",
			"spark", "hdfs", "nginx", "rabbitmq"},
		"spark_pulsar": {"mysql", "python", "fateboard", "client",
			"spark", "hdfs", "nginx", "pulsar"},
		"spark_local_pulsar": {"mysql", "python", "fateboard", "client",
			"nginx", "pulsar"},
	}

	backend, ok := yamlMap["backend"].(string)
	var backendOptions []string
	for k := range backendModules {
		backendOptions = append(backendOptions, k)
	}
	if !ok || !Contains(backend, backendOptions) {
		color.Yellow("Config Warning: the backend in your yaml is not supported\n")
		return
	}

	templateModules := backendModules[backend]
	currentModules, ok := yamlMap["modules"].([]interface{})
	if !ok {
		color.Yellow("Config Warning: the modules in your yaml is not supported\n")
		return
	}

	for _, m := range currentModules {
		// if the module is not in the template modules,
		// alert the user that the module is redundant.
		module := m.(string)
		if !Contains(module, templateModules) {
			color.Yellow("Config Warning: the backend is %s,"+
				" so the redundant module %s is not supported.\n", backend, module)
		}
	}
	for rootKey := range yamlMap {
		// traverse the first layer keys of yaml map to find the modules keys
		// which are in allModules while not in the templateModules.
		if Contains(rootKey, allModules) && !Contains(rootKey, templateModules) {
			color.Yellow("Config Warning: the backend is %s, so whatever configuration "+
				"you have defined about the redundant %s will be ignored.\n", backend, rootKey)
		}
	}
}

func (m *ValidationManager) preprocess() {
	yamlMap := m.testTree.yamlMap
	alertUserIfModulesNotMatchBackend(yamlMap)
}

// ValidateYaml validates the yaml file.
func ValidateYaml(templateValue, testValue string, skippedKeys []string) (errs []error) {
	if templateValue == "" || testValue == "" {
		return []error{errors.New("template or test yaml is empty")}
	}
	templateTree, err := buildValidationTree(templateValue, true, false)
	if err != nil {
		return []error{err}
	}
	testTree, err := buildValidationTree(testValue, false, true)
	if err != nil {
		return []error{err}
	}
	m := &ValidationManager{templateTree, testTree, skippedKeys}

	if m.templateTree == nil || m.testTree == nil {
		return []error{errors.New("building validation tree failed")}
	}
	m.preprocess()
	return m.compareTwoTrees()
}
