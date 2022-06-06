package cli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	root    *TreeNode
	yamlMap map[string]interface{}
	lines   []string
}

type ValidationManager struct {
	templateTree, testTree *ValidationTree
	skippedKeys            []string
}

func trimComments(t []byte) []byte {
	// to trim the comments started with "# "
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

func deconstructKey(k interface{}) (string, int) {
	// to deconstruct the key to the original key and the lineno
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

func NewTreeNode() *TreeNode {
	// return default TreeNode
	node := new(TreeNode)
	node.leaf = false
	node.children = make(map[string]*TreeNode)
	return node
}

func mapToTreeNode(body interface{}, route []string) *TreeNode {
	// recursively convert the yaml map to TreeNode, the route is the path to the current node (if node is a anomymous member in one array, the current route is @ArrayItem). value depends on the type of the key (map, list or basic type)
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

func Contains(element interface{}, set interface{}) bool {
	// check whether an element is in slice/array/map
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

func compareTwoTrees(root_temp, root_test *TreeNode, testLines, skippedKeys []string) (errs []error) {
	// recursively compare the two trees, walking through the nodes not skipped
	valueTemp, valueTest := root_temp.value, root_test.value
	typeTemp, typeTest := reflect.TypeOf(valueTemp), reflect.TypeOf(valueTest)
	if typeTemp != typeTest {
		route := strings.Join(root_test.route, "/")
		errs = append(errs, fmt.Errorf("Your yaml at '%s', line %d \n  '%s' may not match the type\n", route, root_test.lineno, testLines[root_test.lineno]))
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
				errs = append(errs, fmt.Errorf("Your yaml at '%s', line %d \n  '%s' may be redundant\n", route, v.lineno, testLines[v.lineno]))
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

func versionValid(chartVersion string, startVersion []int) (valid bool) {
	// check if the chart version is valid
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

func GetValueTemplateExample(chartName, chartVersion string) (value string, err error) {
	// get the value template example from api
	if !versionValid(chartVersion, []int{1, 9, 0}) {
		err = errors.New("Yaml validation requires the chartVersion >= 1.9.0")
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
		return
	}

	if result, ok := msg.(*ValueResult); ok {
		value = result.Data
	}

	if result, ok := msg.(*ChartResultErr); ok {
		err = errors.New(result.Error)
	}
	return
}

func (m *ValidationManager) compareTwoTrees() []error {
	return compareTwoTrees(m.templateTree.root, m.testTree.root, m.testTree.lines, m.skippedKeys)
}

func stringToReadCloser(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

func readCloserToBuffer(readCloser io.ReadCloser, restoreComments, markLineno bool) ([]byte, []string, error) {
	// read the yaml file and return the content (may be modified) in []byte, the original lines in []string and error
	reader := bufio.NewReader(readCloser)
	defer readCloser.Close()
	buffer := make([]byte, 0, 10)
	lines := make([]string, 1, 10) // 1 for the first lineno

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
			// if one line is pure comment, treat it as a blank line
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

func bufferToMap(buffer []byte) (m map[string]interface{}, err error) {
	// unmarshal the buffer to map, keys are string and numbers are float64
	err = yaml.Unmarshal(buffer, &m)
	return
}

func buildValidationTree(readCloser io.ReadCloser, restoreComments, markLineno bool) (*ValidationTree, error) {
	yamlBuffer, lines, err := readCloserToBuffer(readCloser, restoreComments, markLineno)
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

func getSkippedKeys(m map[string]interface{}) (skippedKeys []string) {
	// the type of skippedKeys array is []interface{}, so we need to transform every element to string
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

func ValidateYaml(templateValue, testValue string, skippedKeys []string) (errs []error) {
	if templateValue == "" || testValue == "" {
		return []error{errors.New("template or test yaml is empty")}
	}
	templateReadCloser, testReadCloser := stringToReadCloser(templateValue), stringToReadCloser(testValue)
	templateTree, err := buildValidationTree(templateReadCloser, true, false)
	if err != nil {
		return []error{err}
	}
	testTree, err := buildValidationTree(testReadCloser, false, true)
	if err != nil {
		return []error{err}
	}
	m := &ValidationManager{templateTree, testTree, skippedKeys}

	if m.templateTree == nil || m.testTree == nil {
		return []error{errors.New("Building validation tree failed")}
	}
	return m.compareTwoTrees()
}
