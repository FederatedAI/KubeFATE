package supportbundle

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	client "github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/k8sclient"
)

type File struct {
	Name string
	Body []byte
}

type Bundler struct {
	KubeConfig string
	CollectDir string
	PackDir    string
	Client     *client.Client
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

// NewBundler returns a bundle helper
func NewBundler(kubeconfig, collectDir, packDir string) (*Bundler, error) {
	c, err := client.NewK8sClient(kubeconfig)
	if err != nil {
		return nil, err
	}
	return &Bundler{kubeconfig, collectDir, packDir, c}, nil
}

// pathExists checks if path exists
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

// CheckDir makes sure one directory is present
func CheckDir(dir string) error {
	if exist, err := pathExists(dir); !exist || err != nil {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// saveFile saves a file
func saveFile(file *File, dir string) error {
	CheckDir(dir)
	f, err := os.Create(path.Join(dir, file.Name))
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(file.Body)
	return nil
}

// dirToZip packs a directory into a zip file
func dirToZip(dir, dst string) error {
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

// CollectIngresses collects ingresses in specific namespace
func (b *Bundler) CollectIngresses(namespace string) error {
	i, err := b.Client.DescribeIngresses(namespace)
	if err != nil {
		return err
	}
	filename := fmt.Sprint("ingresses in ", namespace)
	body := []byte(client.SprintlnIngresses(i))
	file := NewFile(filename, body)
	return saveFile(file, b.CollectDir)
}

// CollectPods collects pods in specific namespace
func (b *Bundler) CollectPods(namespace string, tail int) error {
	p, err := b.Client.DescribePods(namespace, tail)
	if err != nil {
		return err
	}
	filename := fmt.Sprint("pods in ", namespace)
	body := []byte(client.SprintlnPods(p))
	file := NewFile(filename, body)
	return saveFile(file, b.CollectDir)
}

// CollectContainers collects containers in specific namespace
func (b *Bundler) CollectContainers(namespace, podName string, tail int) error {
	s, err := b.Client.DescribePod(namespace, podName, tail)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("containers in %s %s", namespace, podName)
	body := []byte(client.SprintlnContainers(s.Containers))
	file := NewFile(filename, body)
	return saveFile(file, b.CollectDir)
}

// CollectServices collects services in specific namespace
func (b *Bundler) CollectServices(namespace string) error {
	s, err := b.Client.DescribeServices(namespace)
	if err != nil {
		return err
	}
	filename := fmt.Sprint("services in ", namespace)
	body := []byte(client.SprintlnServices(s))
	file := NewFile(filename, body)
	return saveFile(file, b.CollectDir)
}

// CollectNamespace collects ingresses,services,deployments,pods,
// containers as well as their logs in specific namespace and pods into one file
func (b *Bundler) CollectNamespace(namespace string, tail int, podNames ...string) error {
	n, err := b.Client.DescribeNamespace(namespace, tail, podNames...)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	buffer.WriteString(client.SprintlnDeployments(n.Deployments))
	buffer.WriteString(client.SprintlnServices(n.Services))
	buffer.WriteString(client.SprintlnIngresses(n.Ingresses))
	buffer.WriteString(client.SprintlnPods(n.Pods))

	for _, pod := range n.Pods {
		buffer.WriteString(client.SprintlnContainers(pod.Containers))
		for _, c := range pod.Containers {
			filename := fmt.Sprintf("log of container %s in pod %s in namespace %s ",
				c.Name, pod.Name, namespace)
			file := NewFile(filename, []byte(c.Log))
			saveFile(file, b.CollectDir)
		}
	}

	namespaceFilename := fmt.Sprint("namespace ", namespace)
	namespaceFile := NewFile(namespaceFilename, buffer.Bytes())
	err = saveFile(namespaceFile, b.CollectDir)
	return err
}

// Pack packs files into one zip file and removes the CollectDir
func (b *Bundler) Pack() error {
	dst := path.Join(b.PackDir, "supportbundle.zip")
	err := dirToZip(b.CollectDir, dst)
	if err != nil {
		return err
	}
	return os.RemoveAll(b.CollectDir)
}
