package config

import (
	"fmt"

	"github.com/gopi-frame/exception"
)

// ModuleNotFoundException module not found exception
type ModuleNotFoundException struct {
	*exception.Exception
}

// NewModuleNotFoundException new module not found exception
func NewModuleNotFoundException(module string) *ModuleNotFoundException {
	return &ModuleNotFoundException{
		Exception: exception.NewException(fmt.Sprintf("module \"%s\" not found", module)),
	}
}
