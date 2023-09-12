package k8s

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spirt-t/lg-operator/internal/config"
	"github.com/spirt-t/lg-operator/internal/model"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	lgLabelKey                = "kubernetes.generator.label"
	lgPortKey                 = "kubernetes.generator.port"
	creationTimeoutConfigKey  = "kubernetes.timeouts.create"
	deletionTimeoutConfigKey  = "kubernetes.timeouts.delete"
	namespaceKey              = "kubernetes.namespace"
	checkPodReadinessInterval = time.Second * 5
	getExternalIPAttempts     = 5
	getExternalIPInterval     = time.Second * 5
)

//go:generate mockgen -source=./manager.go -destination=./mock/manager.go

// Manager - k8s manager.
type Manager interface {
	Create(ctx context.Context, cfg CreationConfig) (*model.LoadGenerator, error)
	Delete(ctx context.Context, name string) error
	DeleteAll(ctx context.Context) error
	List(ctx context.Context) ([]model.LoadGenerator, error)
}

// CreationConfig for load-generator deploying.
type CreationConfig struct {
	Image            string
	Resources        model.Resources
	Envs             []model.EnvVar
	Commands         []string
	ExposeExternalIP bool
}

type managerImpl struct {
	client    Client
	namespace string
	config    config.Manager
	logger    *zap.Logger
}

// NewManager constructor for Manager.
func NewManager(client Client, config config.Manager, logger *zap.Logger) (Manager, error) {
	var (
		namespace, label string
		err              error
		port             int32
	)

	if er := config.UnmarshalKey(namespaceKey, &namespace); er != nil {
		err = multierr.Append(err, er)
	}

	if er := config.UnmarshalKey(lgLabelKey, &label); er != nil {
		err = multierr.Append(err, er)
	}

	if er := config.UnmarshalKey(lgPortKey, &port); er != nil {
		err = multierr.Append(err, er)
	}

	return &managerImpl{
		client:    client,
		namespace: namespace,
		config:    config,
		logger:    logger,
	}, err
}

func (m *managerImpl) setCreationTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	var timeoutStr string
	if err := m.config.UnmarshalKey(creationTimeoutConfigKey, &timeoutStr); err != nil {
		m.logger.Warn("fail to define creation timeout", zap.Error(err))
		return ctx, func() {}
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		m.logger.Warn("fail to parse creation timeout", zap.Error(err))
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, timeout)
}

// Create new load-generator.
func (m *managerImpl) Create(ctx context.Context, cfg CreationConfig) (*model.LoadGenerator, error) {
	ctx, cancel := m.setCreationTimeout(ctx)
	defer cancel()

	objMeta, err := m.makeObjectMeta()
	if err != nil {
		return nil, err
	}

	var (
		lgService *coreV1.Service
		lgPod     *coreV1.Pod
	)

	defer func() {
		if err != nil {
			// clear k8s resources if failed
			go func(name string) {
				if er := m.Delete(context.Background(), name); er != nil {
					m.logger.Warn("fail to delete k8s entities for generator", zap.Error(er), zap.String("generator_name", name))
				}
			}(objMeta.Name)
		}
	}()

	var port int32
	if err = m.config.UnmarshalKey(lgPortKey, &port); err != nil {
		return nil, fmt.Errorf("fail to define generator port: %w", err)
	}

	_, err = m.createIngress(ctx, objMeta)
	if err != nil {
		return nil, fmt.Errorf("failed to create ingress: %w", err)
	}

	lgService, err = m.createService(ctx, objMeta, port)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	lgPod, err = m.createPod(ctx, cfg, objMeta, port)
	if err != nil {
		return nil, fmt.Errorf("failed to create pod: %w", err)
	}

	var externalIP string
	if cfg.ExposeExternalIP {
		externalIP, err = m.waitServiceExternalIP(ctx, objMeta.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to define service external-ip: %w", err)
		}
	}

	pod := model.LoadGenerator{
		Name:       lgPod.Name,
		ClusterIP:  lgService.Spec.ClusterIP,
		ExternalIP: externalIP,
		Port:       port,
		Status:     lgPod.Status.Phase,
		CreatedAt:  lgPod.CreationTimestamp.Time,
	}

	return &pod, nil
}

func (m *managerImpl) makeObjectMeta() (metaV1.ObjectMeta, error) {
	uuid := uuid.New().String()

	var label string
	if err := m.config.UnmarshalKey(lgLabelKey, &label); err != nil {
		return metaV1.ObjectMeta{}, fmt.Errorf("fail to define label: %w", err)
	}

	return metaV1.ObjectMeta{
		Name: label + "-" + uuid,
		Labels: map[string]string{
			label: "",
		},
	}, nil
}

func (m *managerImpl) createPod(
	ctx context.Context,
	cfg CreationConfig,
	objMeta metaV1.ObjectMeta,
	port int32) (*coreV1.Pod, error) {
	envs := cfg.Envs
	envVars := make([]coreV1.EnvVar, 0, len(envs))

	for _, envVariable := range envs {
		envVars = append(envVars, coreV1.EnvVar{
			Name:  envVariable.Name,
			Value: envVariable.Value,
		})
	}

	resources, err := defineResources(cfg.Resources)
	if err != nil {
		return nil, err
	}

	podConf := coreV1.Pod{
		ObjectMeta: objMeta,
		Spec: coreV1.PodSpec{
			RestartPolicy: coreV1.RestartPolicyNever,
			Containers: []coreV1.Container{{
				Name:            objMeta.Name,
				Image:           cfg.Image,
				ImagePullPolicy: coreV1.PullAlways,
				Env:             envVars,
				Command:         cfg.Commands,
				Resources:       resources,
				Ports: []coreV1.ContainerPort{{
					ContainerPort: port,
				}},
			}},
		},
	}

	lgPod, err := m.client.Get().
		CoreV1().
		Pods(m.namespace).
		Create(ctx, &podConf, metaV1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create pod %s: %w", objMeta.Name, err)
	}

	for {
		select {
		case <-ctx.Done():
			go func(ctx context.Context, name string) {
				if er := m.Delete(ctx, name); er != nil {
					m.logger.Error("fail to delete unready pod "+name, zap.Error(err), zap.String("tank_name", name))
				}
			}(context.Background(), objMeta.Name)

			return nil, fmt.Errorf("context for pod creation exhausted for tank %s; pod will be deleted. Last status: %+v", objMeta.Name, lgPod.Status)
		case <-time.After(checkPodReadinessInterval):
			if lgPod, err = m.client.Get().
				CoreV1().
				Pods(m.namespace).
				Get(ctx, objMeta.Name, metaV1.GetOptions{}); err != nil {
				m.logger.Warn("tank pod status check failed", zap.String("tank_name", objMeta.Name), zap.Error(err))
			}

			if lgPod.Status.Phase == coreV1.PodRunning {
				return lgPod, nil
			}
		}
	}
}

func defineResources(resources model.Resources) (coreV1.ResourceRequirements, error) {
	cpuLimit, err := resource.ParseQuantity(resources.CPU.Limit)
	if err != nil {
		return coreV1.ResourceRequirements{}, fmt.Errorf("fail to parse cpu limit: %w", err)
	}
	cpuRequest, err := resource.ParseQuantity(resources.CPU.Request)
	if err != nil {
		return coreV1.ResourceRequirements{}, fmt.Errorf("fail to parse cpu request: %w", err)
	}
	memoryLimit, err := resource.ParseQuantity(resources.Memory.Limit)
	if err != nil {
		return coreV1.ResourceRequirements{}, fmt.Errorf("fail to parse memory limit: %w", err)
	}
	memoryRequest, err := resource.ParseQuantity(resources.Memory.Request)
	if err != nil {
		return coreV1.ResourceRequirements{}, fmt.Errorf("fail to parse memory request: %w", err)
	}

	return coreV1.ResourceRequirements{
		Limits: coreV1.ResourceList{
			"cpu":    cpuLimit,
			"memory": memoryLimit,
		},
		Requests: coreV1.ResourceList{
			"cpu":    cpuRequest,
			"memory": memoryRequest,
		},
	}, nil
}

func (m *managerImpl) createService(ctx context.Context,
	objMeta metaV1.ObjectMeta,
	port int32) (*coreV1.Service, error) {
	svc, err := m.client.Get().CoreV1().Services(m.namespace).Create(ctx, &coreV1.Service{
		ObjectMeta: objMeta,
		Spec: coreV1.ServiceSpec{
			Type:     coreV1.ServiceTypeLoadBalancer,
			Selector: objMeta.Labels,
			Ports: []coreV1.ServicePort{
				{
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: port,
					},
					Port: port,
				}},
		},
	}, metaV1.CreateOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to create service %s: %w", objMeta.Name, err)
	}

	return svc, nil
}
func (m *managerImpl) waitServiceExternalIP(ctx context.Context, serviceName string) (string, error) {
	var i int

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(getExternalIPInterval):
			if i >= getExternalIPAttempts {
				return "", errors.New("the number of attempts to obtain an external-ip has expired")
			}

			ip, err := m.getServiceExternalIP(ctx, serviceName)
			if err != nil {
				return "", err
			}

			if ip != "" {
				return ip, nil
			}
			i++
		}
	}
}

func (m *managerImpl) getServiceExternalIP(ctx context.Context, serviceName string) (string, error) {
	svc, err := m.client.Get().CoreV1().Services(m.namespace).Get(ctx, serviceName, metaV1.GetOptions{})
	if err != nil {
		return "", err
	}

	var externalIP string
	if len(svc.Status.LoadBalancer.Ingress) > 0 {
		externalIP = svc.Status.LoadBalancer.Ingress[0].IP
	}

	return externalIP, nil
}

func (m *managerImpl) createIngress(ctx context.Context, objMeta metaV1.ObjectMeta) (*v1.Ingress, error) {
	ingress, err := m.client.Get().NetworkingV1().Ingresses(m.namespace).Create(ctx, &v1.Ingress{
		ObjectMeta: objMeta,
		Spec: v1.IngressSpec{
			DefaultBackend: &v1.IngressBackend{
				Resource: &coreV1.TypedLocalObjectReference{
					Kind: "Service",
					Name: objMeta.Name,
				},
			},
		},
	}, metaV1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create ingress %s: %w", objMeta.Name, err)
	}

	return ingress, nil
}

// List of existing load generators.
func (m *managerImpl) List(ctx context.Context) ([]model.LoadGenerator, error) {
	var label string
	if err := m.config.UnmarshalKey(lgLabelKey, &label); err != nil {
		return nil, fmt.Errorf("fail to define label: %w", err)
	}

	podsList, err := m.client.Get().
		CoreV1().
		Pods(m.namespace).
		List(ctx, metaV1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, fmt.Errorf("fail to get list of pods: %w", err)
	}

	servicesList, err := m.client.Get().
		CoreV1().
		Services(m.namespace).
		List(ctx, metaV1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, fmt.Errorf("fail to get list of services: %w", err)
	}

	ipsMp := make(map[string]coreV1.Service)
	for _, service := range servicesList.Items {
		ipsMp[service.Name] = service
	}

	generators := make([]model.LoadGenerator, 0, len(podsList.Items))

	for _, pod := range podsList.Items {
		var externalIP string

		service := ipsMp[pod.Name]
		if len(service.Status.LoadBalancer.Ingress) > 0 {
			externalIP = service.Status.LoadBalancer.Ingress[0].IP
		}

		var port int32
		if len(pod.Spec.Containers) > 0 && len(pod.Spec.Containers[0].Ports) > 0 {
			port = pod.Spec.Containers[0].Ports[0].ContainerPort
		}

		generators = append(generators, model.LoadGenerator{
			Name:       pod.Name,
			ClusterIP:  service.Spec.ClusterIP,
			ExternalIP: externalIP,
			Port:       port,
			Status:     pod.Status.Phase,
			CreatedAt:  pod.CreationTimestamp.Time,
		})
	}

	return generators, nil
}

// Delete load generator by name.
func (m *managerImpl) Delete(ctx context.Context, name string) error {
	ctx, cancel := m.setDeletionTimeout(context.Background())
	defer cancel()

	errPod := m.client.Get().
		CoreV1().
		Pods(m.namespace).
		Delete(ctx, name, metaV1.DeleteOptions{})

	errService := m.client.Get().
		CoreV1().
		Services(m.namespace).
		Delete(ctx, name, metaV1.DeleteOptions{})

	errIngress := m.client.Get().
		NetworkingV1().
		Ingresses(m.namespace).
		Delete(ctx, name, metaV1.DeleteOptions{})

	return multierr.Combine(errPod, errService, errIngress)
}

// DeleteAll generators.
func (m *managerImpl) DeleteAll(ctx context.Context) error {
	ctx, cancel := m.setDeletionTimeout(context.Background())
	defer cancel()

	var label string
	var err error
	if err = m.config.UnmarshalKey(lgLabelKey, &label); err != nil {
		return fmt.Errorf("fail to define label: %w", err)
	}

	errPod := m.client.Get().
		CoreV1().
		Pods(m.namespace).
		DeleteCollection(ctx, metaV1.DeleteOptions{}, metaV1.ListOptions{
			LabelSelector: label,
		})

	servicesList, errService := m.client.Get().
		CoreV1().
		Services(m.namespace).
		List(ctx, metaV1.ListOptions{LabelSelector: label})
	if errService == nil {
		for _, svc := range servicesList.Items {
			if er := m.client.Get().
				CoreV1().
				Services(m.namespace).
				Delete(ctx, svc.Name, metaV1.DeleteOptions{}); er != nil {
				errService = multierr.Append(errService, er)
			}
		}
	} else {
		errService = fmt.Errorf("fail to get list of services: %w", err)
	}

	errIngress := m.client.Get().
		NetworkingV1().
		Ingresses(m.namespace).
		DeleteCollection(ctx, metaV1.DeleteOptions{}, metaV1.ListOptions{
			LabelSelector: label,
		})

	return multierr.Combine(errPod, errService, errIngress)
}

func (m *managerImpl) setDeletionTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	var timeoutStr string
	if err := m.config.UnmarshalKey(deletionTimeoutConfigKey, &timeoutStr); err != nil {
		m.logger.Warn("fail to define deletion timeout", zap.Error(err))
		return ctx, func() {}
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		m.logger.Warn("fail to parse deletion timeout", zap.Error(err))
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, timeout)
}
