package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Manager for user config.
type Manager interface {
	UnmarshalKey(key string, val interface{}) error
}

type managerImpl struct {
	v *viper.Viper
}

// NewManager - constructor for Manager.
func NewManager(cfgPath string) (Manager, error) {
	v := viper.New()
	v.SetConfigFile(cfgPath)
	v.SetConfigType(configType(cfgPath))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &managerImpl{v}, nil
}

func configType(cfgPath string) string {
	parts := strings.Split(cfgPath, ".")

	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}

// UnmarshalKey - read config value by key.
func (m *managerImpl) UnmarshalKey(key string, val interface{}) error {
	return m.v.UnmarshalKey(key, val)
}
