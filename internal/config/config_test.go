package config

import (
	"testing"

	"github.com/spirt-t/lg-operator/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	mngr, err := NewManager("../../testfiles/test_config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("service ports", func(t *testing.T) {
		var httpPort, grpcPort int
		err = mngr.UnmarshalKey("service.ports.http", &httpPort)
		assert.NoError(t, err)
		assert.Equal(t, 7000, httpPort)

		err = mngr.UnmarshalKey("service.ports.grpc", &grpcPort)
		assert.NoError(t, err)
		assert.Equal(t, 7002, grpcPort)
	})

	t.Run("kubernetes", func(t *testing.T) {
		var host, namespace, label string
		var port int

		err = mngr.UnmarshalKey("kubernetes.service.host", &host)
		assert.NoError(t, err)
		assert.Equal(t, "10.96.0.1", host)

		err = mngr.UnmarshalKey("kubernetes.service.port", &port)
		assert.NoError(t, err)
		assert.Equal(t, 443, port)

		err = mngr.UnmarshalKey("kubernetes.namespace", &namespace)
		assert.NoError(t, err)
		assert.Equal(t, "default", namespace)

		err = mngr.UnmarshalKey("kubernetes.generator.label", &label)
		assert.NoError(t, err)
		assert.Equal(t, "load-generator", label)
	})

	t.Run("resources model", func(t *testing.T) {
		var resources model.Resources

		err = mngr.UnmarshalKey("default_resources", &resources)
		assert.NoError(t, err)
		assert.Equal(t, "1", resources.CPU.Request)
		assert.Equal(t, "2", resources.CPU.Limit)
		assert.Equal(t, "1Gi", resources.Memory.Request)
		assert.Equal(t, "2Gi", resources.Memory.Limit)
	})

	t.Run("cleaner enabled", func(t *testing.T) {
		var enabled bool
		err = mngr.UnmarshalKey("cleaning.completed.enabled", &enabled)
		assert.NoError(t, err)
		assert.True(t, enabled)
	})
}
