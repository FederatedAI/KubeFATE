package k8sclient

import (
	"context"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/gosuri/uitable"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	*kubernetes.Clientset
}

type Namespace struct {
	Name        string
	Ingresses   []*Ingress
	Services    []*Service
	Deployments []*Deployment
	Pods        []*Pod
}

type Container struct {
	Name        string
	Status      string
	Ready       bool
	Image       string
	ImageID     string
	ContainerID string
	Log         string
}

type Pod struct {
	Name              string
	Status            string
	Ready             string
	Containers        []*Container
	CreationTimestamp time.Time
}

type Deployment struct {
	Name                string
	Replicas            int
	UpdatedReplicas     int
	ReadyReplicas       int
	AvailableReplicas   int
	UnavailableReplicas int
}

type Service struct {
	Name              string
	Type              string
	ClusterIP         string
	ExternalIPs       []string
	Ports             []string
	CreationTimestamp time.Time
}

type Ingress struct {
	Name              string
	Class             string
	Host              string
	Address           string
	Ports             []string
	CreationTimestamp time.Time
}

// NewK8sClient returns a k8s restClient
func NewK8sClient(kubeconfig string) (c *Client, err error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return
	}
	client, err := kubernetes.NewForConfig(config)
	c = &Client{client}
	return
}

// NamespaceListToNames extracts and returns namespace list
func NamespaceListToNames(list *v1.NamespaceList) []string {
	var names []string
	for _, i := range list.Items {
		names = append(names, i.Name)
	}
	return names
}

// PodListToNames extracts and returns podname list
func PodListToNames(list *v1.PodList) []string {
	var names []string
	for _, i := range list.Items {
		names = append(names, i.Name)
	}
	return names
}

// GetNamespaceList returns namespaceList from k8s api
func (c *Client) GetNamespaceList() (*v1.NamespaceList, error) {
	return c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
}

// GetNamespaceList returns specific namespace from k8s api
func (c *Client) GetNamespace(namespace string) (*v1.Namespace, error) {
	return c.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
}

// ToNamespace tansforms v1.Namespace to Namespace
func ToNamespace(n *v1.Namespace) *Namespace {
	return &Namespace{
		Name: n.Name,
	}
}

// DescribeNamespace gets and returns all msg in specific namespace
func (c *Client) DescribeNamespace(
	namespace string, tail int, podNames ...string) (n *Namespace, err error) {
	n = &Namespace{Name: namespace}
	d, err := c.DescribeDeployments(namespace)
	if err != nil {
		return
	}
	n.Deployments = d

	i, err := c.DescribeIngresses(namespace)
	if err != nil {
		return
	}
	n.Ingresses = i

	s, err := c.DescribeServices(namespace)
	if err != nil {
		return
	}
	n.Services = s

	var pods []*Pod
	for _, pod := range podNames {
		p, err := c.DescribePod(namespace, pod, tail)
		if err != nil {
			continue
		}
		pods = append(pods, p)
	}
	n.Pods = pods
	return
}

// GetDeploymentList returns DeploymentList from k8s api
func (c *Client) GetDeploymentList(namespace string) (*appsv1.DeploymentList, error) {
	return c.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
}

// ToDeployment tansforms appsv1.Deployment to Deployment
func ToDeployment(d *appsv1.Deployment) *Deployment {
	return &Deployment{
		Name:                d.Name,
		Replicas:            int(d.Status.Replicas),
		UpdatedReplicas:     int(d.Status.UpdatedReplicas),
		ReadyReplicas:       int(d.Status.ReadyReplicas),
		AvailableReplicas:   int(d.Status.AvailableReplicas),
		UnavailableReplicas: int(d.Status.UnavailableReplicas),
	}
}

// DescribeDeployments returns DeploymentList
func (c *Client) DescribeDeployments(namespace string) ([]*Deployment, error) {
	d, err := c.GetDeploymentList(namespace)
	if err != nil {
		return nil, err
	}
	var deployments []*Deployment
	for _, i := range d.Items {
		deployments = append(deployments, ToDeployment(&i))
	}
	return deployments, nil
}

// SprintlnDeployments returns tabular output of DeploymentList
func SprintlnDeployments(items []*Deployment) string {
	table := uitable.New()
	table.AddRow("Deployments")
	table.AddRow("Name", "Replicas", "UpdatedReplicas",
		"ReadyReplicas", "AvailableReplicas", "UnavailableReplicas")
	for _, r := range items {
		table.AddRow(r.Name, r.Replicas, r.UpdatedReplicas,
			r.ReadyReplicas, r.AvailableReplicas, r.UnavailableReplicas)
	}
	table.AddRow("")
	return fmt.Sprintln(table)
}

func (c *Client) GetServiceList(namespace string) (*v1.ServiceList, error) {
	return c.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetService(namespace, service string) (*v1.Service, error) {
	return c.CoreV1().Services(namespace).Get(context.TODO(), service, metav1.GetOptions{})
}

// ToPorts tansforms v1.ServicePort to ports string
func ToPorts(p []v1.ServicePort) []string {
	var ports []string
	for _, i := range p {
		ports = append(ports, fmt.Sprintf("%d/%s", i.Port, i.Protocol))
	}
	if len(ports) == 0 {
		ports = append(ports, "<nil>")
	}
	return ports
}

// ToService tansforms v1.Service to Service
func ToService(s *v1.Service) *Service {
	return &Service{
		Name:              s.Name,
		Type:              string(s.Spec.Type),
		ClusterIP:         s.Spec.ClusterIP,
		ExternalIPs:       s.Spec.ExternalIPs,
		Ports:             ToPorts(s.Spec.Ports),
		CreationTimestamp: s.CreationTimestamp.Time,
	}
}

// DescribeServices returns ServiceList
func (c *Client) DescribeServices(namespace string) ([]*Service, error) {
	s, err := c.GetServiceList(namespace)
	if err != nil {
		return nil, err
	}
	var services []*Service
	for _, i := range s.Items {
		services = append(services, ToService(&i))
	}
	return services, nil
}

// SprintlnDeployments returns tabular output of DeploymentList
func SprintlnServices(items []*Service) string {
	table := uitable.New()
	table.AddRow("Services")
	table.AddRow("Name", "Type", "ClusterIP", "ExternalIPs", "Ports", "CreationTime")
	for _, r := range items {
		table.AddRow(r.Name, r.Type, r.ClusterIP,
			r.ExternalIPs, r.CreationTimestamp.Format(time.RFC3339))
	}
	table.AddRow("")
	return fmt.Sprintln(table)
}

// GetIngressList returns IngressList from k8s api
func (c *Client) GetIngressList(namespace string) (*networkingv1.IngressList, error) {
	return c.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
}

// ToIngress tansforms networkingv1.Ingress to Ingress
func ToIngress(i *networkingv1.Ingress) *Ingress {
	return &Ingress{
		Name:              i.Name,
		Class:             *i.Spec.IngressClassName,
		Host:              i.Spec.Rules[0].Host,
		Address:           i.Status.LoadBalancer.Ingress[0].IP,
		CreationTimestamp: i.CreationTimestamp.Time,
	}
}

// DescribeIngresses returns IngressList
func (c *Client) DescribeIngresses(namespace string) ([]*Ingress, error) {
	i, err := c.GetIngressList(namespace)
	if err != nil {
		return nil, err
	}
	var ingresses []*Ingress
	for _, ingress := range i.Items {
		ingresses = append(ingresses, ToIngress(&ingress))
	}
	return ingresses, nil
}

// SprintlnIngresses returns tabular output of IngressList
func SprintlnIngresses(items []*Ingress) string {
	table := uitable.New()
	table.AddRow("Ingresses")
	table.AddRow("Name", "Class", "Host", "Address", "Ports", "CreationTime")
	for _, r := range items {
		table.AddRow(r.Name, r.Class, r.Host, r.Ports, r.CreationTimestamp.Format(time.RFC3339))
	}
	table.AddRow("")
	return fmt.Sprintln(table)
}

// GetPodList returns PodList from k8s api
func (c *Client) GetPodList(namespace string) (*v1.PodList, error) {
	return c.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}

// GetPod returns specific Pod from k8s api
func (c *Client) GetPod(namespace, pod string) (*v1.Pod, error) {
	return c.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
}

// ToPod tansforms v1.Pod to Pod
func ToPod(p *v1.Pod) *Pod {
	return &Pod{
		Name:              p.Name,
		Status:            string(p.Status.Phase),
		CreationTimestamp: p.CreationTimestamp.Time,
	}
}

// GetContainersFromPod fills a pod with it's containers
func (c *Client) GetContainersFromPod(p *v1.Pod, tail int) ([]*Container, error) {
	var containers []*Container
	for i := range p.Spec.Containers {
		container := p.Spec.Containers[i]
		status := p.Status.ContainerStatuses[i]
		log, err := c.GetLog(p.Namespace, p.Name, container.Name, tail)
		if err != nil {
			continue
		}
		_container := &Container{
			Name:        container.Name,
			Status:      status.State.String(),
			Ready:       status.Ready,
			Image:       container.Image,
			ImageID:     status.ImageID,
			ContainerID: status.ContainerID,
			Log:         log,
		}
		containers = append(containers, _container)
	}
	return containers, nil
}

// DescribePod returns specific Pod
func (c *Client) DescribePod(namespace, podName string, tail int) (*Pod, error) {
	p, err := c.GetPod(namespace, podName)
	if err != nil {
		return nil, err
	}
	pod := ToPod(p)
	pod.Containers, _ = c.GetContainersFromPod(p, tail)
	return pod, nil
}

// DescribePod returns PodList
func (c *Client) DescribePods(namespace string, tail int) ([]*Pod, error) {
	p, err := c.GetPodList(namespace)
	if err != nil {
		return nil, err
	}
	var pods []*Pod
	for _, i := range p.Items {
		if pod, err := c.DescribePod(namespace, i.Name, tail); err == nil {
			pods = append(pods, pod)
		}
	}
	return pods, nil
}

// SprintlnPods returns tabular output of PodList
func SprintlnPods(items []*Pod) string {
	table := uitable.New()
	table.AddRow("Pods")
	table.AddRow("Name", "Status", "Ready", "Containers", "CreationTime")
	for _, r := range items {
		var names []string
		for _, c := range r.Containers {
			names = append(names, c.Name)
		}
		table.AddRow(r.Name, r.Status, r.Ready, names, r.CreationTimestamp.Format(time.RFC3339))
	}
	table.AddRow("")
	return fmt.Sprintln(table)
}

// SprintlnContainers returns tabular output of ContainerList
func SprintlnContainers(items []*Container) string {
	table := uitable.New()
	table.AddRow("Containers")
	table.AddRow("Name", "Status", "Ready", "Image", "ImageID", "ContainerID")
	for _, r := range items {
		table.AddRow(r.Name, r.Status, r.Ready, r.Image, r.ImageID, r.ContainerID)
	}
	table.AddRow("")
	return fmt.Sprintln(table)
}

// GetLog returns log of one container
func (c *Client) GetLog(namespace, pod, container string, tail int) (string, error) {
	tailLine := int64(tail)
	logOpt := &v1.PodLogOptions{
		Container: container,
		TailLines: &tailLine,
	}
	req := c.CoreV1().Pods(namespace).GetLogs(pod, logOpt)
	r, err := req.Stream(context.TODO())
	if err != nil {
		return "", err
	}
	defer r.Close()
	buffer, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

// GetKubeconfig gets default kubeconfig
func GetKubeconfig() string {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "",
			"absolute path to the kubeconfig file")
	}
	flag.Parse()
	return *kubeconfig
}
