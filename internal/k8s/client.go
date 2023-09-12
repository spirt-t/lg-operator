package k8s

import (
	"context"
	"fmt"
	"os"

	"github.com/spirt-t/lg-operator/internal/config"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	k8sHostKey = "kubernetes.service.host"
	k8sPortKey = "kubernetes.service.port"
)

// Client to k8s.
type Client interface {
	Init(ctx context.Context) error
	Get() *kubernetes.Clientset
}

// NewClient constructor for k8s client.
func NewClient(config config.Manager, logger *zap.Logger) Client {
	return &clientImpl{config: config, logger: logger}
}

type clientImpl struct {
	client *kubernetes.Clientset
	config config.Manager
	logger *zap.Logger
}

// Init client.
func (c *clientImpl) Init(_ context.Context) error {
	var (
		host, port string
	)

	err := c.config.UnmarshalKey(k8sHostKey, &host)
	if err != nil {
		return fmt.Errorf("fail to get parameter %s: %w", k8sHostKey, err)
	}

	err = c.config.UnmarshalKey(k8sPortKey, &port)
	if err != nil {
		return fmt.Errorf("fail to get parameter %s: %w", k8sPortKey, err)
	}

	os.Setenv("KUBERNETES_SERVICE_HOST", host)
	os.Setenv("KUBERNETES_SERVICE_PORT", port)

	configKuber, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("fail to read k8s cluster config: %w", err)
	}

	c.client, err = kubernetes.NewForConfig(configKuber)
	if err != nil {
		return fmt.Errorf("fail to make k8s client: %w", err)
	}

	info, err := c.client.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("fail to make k8s client: %w", err)
	}

	c.logger.Info("k8s server version", zap.Stringer("info", info))

	return nil
}

// Get connection to k8s.
func (c *clientImpl) Get() *kubernetes.Clientset {
	return c.client
}
