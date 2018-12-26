package stub

import (
	"fmt"
	"github.com/Percona-Lab/percona-server-mongodb-operator/internal/mongod"
	"github.com/Percona-Lab/percona-server-mongodb-operator/internal/sdk"
	"github.com/Percona-Lab/percona-server-mongodb-operator/internal/util"
	"github.com/Percona-Lab/percona-server-mongodb-operator/pkg/apis/psmdb/v1alpha1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"time"
)

func (h *Handler) ensureExtServices(m *v1alpha1.PerconaServerMongoDB, replset *v1alpha1.ReplsetSpec, podList *corev1.PodList) ([]corev1.Service, error) {
	services := make([]corev1.Service, 0)
	for _, pod := range podList.Items {
		logrus.Infof("Checking that pod %s of replset %s has attached service", pod.Name, replset.Name)

		meta := serviceMeta(m.Namespace, pod.Name)

		logrus.Debugf("Service meta: %v", meta)

		if err := h.client.Get(meta); err != nil {
			if errors.IsNotFound(err) {
				logrus.Infof("pod %s of replset %s doesn't have attached service", pod.Name, replset.Name)
				svc := extService(m, pod.Name)
				if err := createExtService(h.client, svc); err != nil {
					return nil, fmt.Errorf("failed to create external service for replset %s: %v", replset.Name, err)
				}
				meta = svc
			} else {
				return nil, fmt.Errorf("failed to fetch service for replset %s: %v", replset.Name, err)
			}
		}

		logrus.Infof("service %s for pod %s of repleset %s has been found", meta.Name, pod.Name, replset.Name)

		if err := updateExtService(h.client, meta); err != nil {
			return nil, fmt.Errorf("failed to update external service for replset %s: %v", replset.Name, err)
		}
		services = append(services, *meta)
	}

	return services, nil
}

func getExtServices(m *v1alpha1.PerconaServerMongoDB, podName string) (*corev1.Service, error) {
	client := sdk.NewClient()
	svcMeta := serviceMeta(m.Namespace, podName)
	if err := client.Get(svcMeta); err != nil {
		if !errors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to fetch service: %v", err)
		}
	}
	return svcMeta, nil
}

func createExtService(cli sdk.Client, svc *corev1.Service) error {
	logrus.Infof("Creating service %s", svc.Name)
	if err := cli.Create(svc); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		logrus.Infof("Service %s already exist. Skipping", svc.Name)
	}
	return nil
}

func updateExtService(cli sdk.Client, svc *corev1.Service) error {
	var retries uint64 = 0

	for retries <= 5 {
		if err := cli.Update(svc); err != nil {
			if errors.IsConflict(err) {
				time.Sleep(500 * time.Millisecond)
				retries += 1
				continue
			} else {
				return fmt.Errorf("failed to update service: %v", err)
			}
		}
		return nil
	}
	return fmt.Errorf("failed to update service %s, retries limit reached", svc.Name)
}

func serviceMeta(namespace, podName string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
	}
}

func extService(m *v1alpha1.PerconaServerMongoDB, podName string) *corev1.Service {
	svc := serviceMeta(m.Namespace, podName)
	svc.Spec = corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       mongod.MongodPortName,
				Port:       m.Spec.Mongod.Net.Port,
				TargetPort: intstr.FromInt(int(m.Spec.Mongod.Net.Port)),
			},
		},
		Selector: map[string]string{"statefulset.kubernetes.io/pod-name": podName},
	}
	switch m.Spec.Expose.ExposeType {
	case corev1.ServiceTypeNodePort:
		svc.Spec.Type = corev1.ServiceTypeNodePort
		svc.Spec.ExternalTrafficPolicy = "Local"
	case corev1.ServiceTypeLoadBalancer:
		svc.Spec.Type = corev1.ServiceTypeLoadBalancer
		svc.Spec.ExternalTrafficPolicy = "Local"
	default:
		svc.Spec.Type = corev1.ServiceTypeClusterIP
	}
	util.AddOwnerRefToObject(svc, util.AsOwner(m))
	return svc
}

func getServiceAddr(svc corev1.Service, pod corev1.Pod) string {
	var hostname string
	var port int

	switch svc.Spec.Type {
	case corev1.ServiceTypeClusterIP:
		hostname = svc.Spec.ClusterIP
		for _, p := range svc.Spec.Ports {
			if p.Name != mongod.MongodPortName {
				continue
			}
			port = int(p.Port)
		}

	case corev1.ServiceTypeLoadBalancer:
		hostname = svc.Spec.LoadBalancerIP
		for _, p := range svc.Spec.Ports {
			if p.Name != mongod.MongodPortName {
				continue
			}
			port = int(p.Port)
		}

	case corev1.ServiceTypeNodePort:
		hostname = pod.Status.HostIP
		for _, p := range svc.Spec.Ports {
			if p.Name != mongod.MongodPortName {
				continue
			}
			port = int(p.NodePort)
		}
	}
	return hostname + ":" + strconv.Itoa(port)
}
