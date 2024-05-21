package config

import (
	"github.com/gopi-frame/contract/container"
	"github.com/gopi-frame/contract/support"
)

// ServerProvider server provider
type ServerProvider struct {
	support.ServerProvider
}

// Register register
func (s *ServerProvider) Register(c container.Container) {
	c.Bind("config", func(c container.Container) any {
		return NewRepository()
	})
}
