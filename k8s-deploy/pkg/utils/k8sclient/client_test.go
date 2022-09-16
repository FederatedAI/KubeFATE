package k8sclient

import (
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	namespace1 = v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "1"}}
	namespace2 = v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "2"}}

	deployment1 = appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "1"}}
	deployment2 = appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "2"}}

	servicePorts1 = []v1.ServicePort{{Protocol: "tcp", Port: 8080}, {Protocol: "udp", Port: 8081}}
	servicePorts2 = []v1.ServicePort{}

	service1 = v1.Service{ObjectMeta: metav1.ObjectMeta{Name: "1"}, Spec: v1.ServiceSpec{Ports: servicePorts1}}
	service2 = v1.Service{ObjectMeta: metav1.ObjectMeta{Name: "2"}, Spec: v1.ServiceSpec{Ports: servicePorts2}}

	pathType1        = networkingv1.PathType("pathType1")
	pathType2        = networkingv1.PathType("pathType2")
	paths            = []networkingv1.HTTPIngressPath{{Path: "/1", PathType: &pathType1}, {Path: "/2", PathType: &pathType2}}
	ingressRules     = []networkingv1.IngressRule{{Host: "example.com", IngressRuleValue: networkingv1.IngressRuleValue{HTTP: &networkingv1.HTTPIngressRuleValue{Paths: paths}}}}
	ingressClassName = "cls1"
	ingress1         = networkingv1.Ingress{ObjectMeta: metav1.ObjectMeta{Name: "1"}, Spec: networkingv1.IngressSpec{IngressClassName: &ingressClassName, Rules: ingressRules}}

	pod1 = v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "1"}}
	pod2 = v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "2"}}
)

func TestNewK8sClient(t *testing.T) {
	type args struct {
		kubeconfig string
	}
	tests := []struct {
		name    string
		args    args
		wantC   *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewK8sClient(tt.args.kubeconfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewK8sClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewK8sClient() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestNamespaceListToNames(t *testing.T) {
	type args struct {
		list *v1.NamespaceList
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{"test", args{&v1.NamespaceList{Items: []v1.Namespace{namespace1, namespace2}}}, []string{"1", "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NamespaceListToNames(tt.args.list); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NamespaceListToNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodListToNames(t *testing.T) {
	type args struct {
		list *v1.PodList
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"test", args{&v1.PodList{Items: []v1.Pod{pod1, pod2}}}, []string{"1", "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PodListToNames(tt.args.list); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodListToNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNamespace(t *testing.T) {
	type args struct {
		n *v1.Namespace
	}
	tests := []struct {
		name string
		args args
		want *Namespace
	}{
		{"1", args{&namespace1}, &Namespace{Name: "1"}},
		{"2", args{&namespace2}, &Namespace{Name: "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToNamespace(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDeployment(t *testing.T) {
	type args struct {
		d *appsv1.Deployment
	}
	tests := []struct {
		name string
		args args
		want *Deployment
	}{
		{"1", args{&deployment1}, &Deployment{Name: "1"}},
		{"2", args{&deployment2}, &Deployment{Name: "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToDeployment(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPorts(t *testing.T) {
	type args struct {
		p []v1.ServicePort
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"notEmpty", args{servicePorts1}, []string{"8080/tcp", "8081/udp"}},
		{"empty", args{servicePorts2}, []string{"<nil>"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPorts(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPorts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToService(t *testing.T) {
	type args struct {
		s *v1.Service
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{"1", args{&service1}, &Service{Name: "1", Ports: []string{"8080/tcp", "8081/udp"}}},
		{"1", args{&service2}, &Service{Name: "2", Ports: []string{"<nil>"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToService(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToIngress(t *testing.T) {
	type args struct {
		i *networkingv1.Ingress
	}
	tests := []struct {
		name string
		args args
		want *Ingress
	}{
		{"1", args{&ingress1}, &Ingress{Name: "1", Class: "cls1", Rules: []IngressRule{{Host: "example.com", Path: []HostPath{{Path: "/1", PathType: string(pathType1), Backend: ""}, {Path: "/2", PathType: string(pathType2), Backend: ""}}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToIngress(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToIngress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPod(t *testing.T) {
	type args struct {
		p *v1.Pod
	}
	tests := []struct {
		name string
		args args
		want *Pod
	}{
		{"1", args{&pod1}, &Pod{Name: "1"}},
		{"2", args{&pod2}, &Pod{Name: "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPod(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToIngressRules(t *testing.T) {
	type args struct {
		rules []networkingv1.IngressRule
	}
	tests := []struct {
		name string
		args args
		want []IngressRule
	}{
		{"1", args{ingressRules}, []IngressRule{{Host: "example.com", Path: []HostPath{{Path: "/1", PathType: string(pathType1), Backend: ""}, {Path: "/2", PathType: string(pathType2), Backend: ""}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToIngressRules(tt.args.rules); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToIngressRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
