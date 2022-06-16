/*
 * Copyright 2019-2021 VMware, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package cli

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

const (
	bundleDir         = "./"
	temporaryFilesDir = "./kubefate-supportbundle"
)

func SupportBundleCommand() *cli.Command {
	return &cli.Command{
		Name:  "supportbundle",
		Flags: []cli.Flag{},
		Subcommands: []*cli.Command{
			SupportBundleCollectCommand(),
			SupportBundlePackCommand(),
		},
		Usage: "Collect data to perform troubleshooting",
	}
}

func SupportBundleCollectCommand() *cli.Command {
	return &cli.Command{
		Name: "collect",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "tail",
				Aliases: []string{"t"},
				Value:   200,
				Usage:   "Specify how many rows to record.",
			},
		},
		Usage: "Collect data",
		Action: func(c *cli.Context) error {
			tail := c.Int("tail")
			return Collect(tail)
		},
	}
}

func SupportBundlePackCommand() *cli.Command {
	return &cli.Command{
		Name:  "pack",
		Flags: []cli.Flag{},
		Usage: "Pack data",
		Action: func(c *cli.Context) error {
			return Pack()
		},
	}
}

// RunCommand runs command and returns output
func RunCommand(cmd string) ([]byte, error) {
	c := exec.Command("bash", "-c", cmd)
	return c.CombinedOutput()
}

type File struct {
	Name string
	Body []byte
}

// NewFile creates struct File with space-trimed and dash-joined name
func NewFile(name string, body []byte) *File {
	group := strings.Split(name, " ")
	newGroup := make([]string, 0, len(group))
	for _, s := range group {
		if s != "" {
			newGroup = append(newGroup, s)
		}
	}
	newName := strings.Join(newGroup, "-")
	return &File{newName, body}
}

func kubectlLogs(namespace, pod string, tail int) ([]byte, string, error) {
	cmd := fmt.Sprintf("kubectl logs --tail=%d -n %s %s ", tail, namespace, pod)
	log, err := RunCommand(cmd)
	return log, cmd, err
}

// getShellOutputColumns returns values with the given column name
func getShellOutputColumns(cmd, column string) ([][]byte, error) {
	cmdFirstRow := fmt.Sprintf("%s | awk '{if (NR==1) print}'", cmd)
	outputFirstRow, err := RunCommand(cmdFirstRow)
	if err != nil {
		return nil, err
	}
	columnIdx := -1
	for idx, c := range bytes.Split(outputFirstRow, []byte(" ")) {
		if string(c) == column {
			columnIdx = idx
		}
	}
	if columnIdx < 0 {
		return nil, errors.New("column not found")
	}
	cmdColumn := fmt.Sprintf("%s | awk '{if (NR>1) {print $%d}}'", cmd, columnIdx+1)
	outputColumn, err := RunCommand(cmdColumn)
	if err != nil {
		return nil, err
	}
	return bytes.Split(bytes.Trim(outputColumn, "\n"), []byte("\n")), nil
}

func collectPods() (files []*File, err error) {
	cmdPods := "kubectl get pods -A"
	podsNames, err := getShellOutputColumns(cmdPods, "NAME")
	if err != nil {
		return
	}

	var names []string
	for _, name := range podsNames {
		names = append(names, string(name))
	}

	pods := []string{}
	podsPrompt := &survey.MultiSelect{
		Message: "What pods do you want to share logs with us:",
		Options: names,
	}
	survey.AskOne(podsPrompt, &pods)

	for _, pod := range pods {
		log, cmd, err := kubectlLogs("kube-fate", string(pod), 200)
		if err != nil {
			return nil, err
		}
		files = append(files, NewFile(cmd, log))
	}
	return
}

// collectKubeFATE returns files related to KubeFATE
func collectKubeFATE(tail int) (files []*File, err error) {
	cmdAll := "kubectl get all -n kube-fate"
	outputAll, err := RunCommand(cmdAll)
	if err != nil {
		return
	}
	files = append(files, NewFile(cmdAll, outputAll))

	cmdPods := "kubectl get pods  -n kube-fate"
	podsNames, err := getShellOutputColumns(cmdPods, "NAME")
	if err != nil {
		return
	}
	for _, pod := range podsNames {
		log, cmd, err := kubectlLogs("kube-fate", string(pod), tail)
		if err != nil {
			return nil, err
		}
		files = append(files, NewFile(cmd, log))
	}
	return
}

// saveFile saves a file
func saveFile(file *File, dir string) error {
	f, err := os.Create(path.Join(dir, file.Name))
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(file.Body)
	return nil
}

// zipDir packs a directory into a zip file
func zipDir(dir, dst string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	fileDst, _ := os.Create(dst)
	w := zip.NewWriter(fileDst)
	defer w.Close()
	for _, file := range files {
		fw, err := w.Create(file.Name())
		if err != nil {
			return err
		}
		fileContent, err := ioutil.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			return err
		}
		_, err = fw.Write(fileContent)
		if err != nil {
			return err
		}
	}
	return nil
}

// pathExists check if path exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Collect gathers and saves information into files
func Collect(tail int) error {
	var files []*File

	deployMethod := ""
	deployPrompt := &survey.Select{
		Message: "Choose way of deploy:",
		Options: []string{"k8s", "docker"},
		Default: "k8s",
	}
	survey.AskOne(deployPrompt, &deployMethod)

	if f, err := collectKubeFATE(tail); err != nil {
		fmt.Println(err)
		return err
	} else {
		files = append(files, f...)
	}

	if exist, _ := pathExists(temporaryFilesDir); !exist {
		os.MkdirAll(temporaryFilesDir, os.ModePerm)
	}
	for _, file := range files {
		if err := saveFile(file, temporaryFilesDir); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("You can attach additional files about this issue into %s\n", temporaryFilesDir)
	return nil
}

// Pack packs files into one zip file and removes the temporary directory
func Pack() error {
	privacyAck := true
	privacyPrompt := &survey.Confirm{
		Message: "This program collects some data, please ensure that there is no any privacy issue.\n Enter y if you're willing to share these collected data with us.",
	}
	survey.AskOne(privacyPrompt, &privacyAck)
	if !privacyAck {
		return errors.New("Please acknowledge the privacy agreement")
	}
	dst := path.Join(bundleDir, "supportbundle.zip")
	err := zipDir(temporaryFilesDir, dst)
	if err != nil {
		return err
	}
	err = os.RemoveAll(temporaryFilesDir)
	return err
}
