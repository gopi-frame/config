package config

import (
	"github.com/gopi-frame/contract/container"
)

// ServerProvider server provider
type ServerProvider struct{}

// Register register
func (s *ServerProvider) Register(c container.Container) {
	c.Bind("config", func(c container.Container) any {
		return NewRepository()
	})
}
