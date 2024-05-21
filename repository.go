package config

import (
	"io"

	"github.com/gopi-frame/support/maps"
	"github.com/spf13/viper"
)

// NewRepository new repository
func NewRepository() *Repository {
	r := new(Repository)
	r.modules = maps.NewMap[string, *viper.Viper]()
	return r
}

// Repository repository
type Repository struct {
	modules *maps.Map[string, *viper.Viper]
}

// Module get module
func (repo *Repository) Module(module string) *viper.Viper {
	if repo.modules.TryLock() {
		defer repo.modules.Unlock()
	}
	return repo.modules.GetOr(module, nil)
}

func (repo *Repository) Read(module string, reader io.Reader) error {
	if repo.modules.TryLock() {
		defer repo.modules.Unlock()
	}
	viper := viper.New()
	err := viper.ReadConfig(reader)
	if err != nil {
		return err
	}
	repo.modules.Set(module, viper)
	return nil
}

// Set set
func (repo *Repository) Set(module string, key string, value any) {
	if repo.modules.TryLock() {
		defer repo.modules.Unlock()
	}
	repo.Module(module).Set(key, value)
}

// Has has
func (repo *Repository) Has(module string, key string) bool {
	if repo.modules.TryLock() {
		defer repo.modules.Unlock()
	}
	if module := repo.Module(module); module != nil {
		return module.IsSet(key)
	}
	return false
}

// Get get
func (repo *Repository) Get(module string, key string, defaultValue ...any) any {
	if repo.modules.TryLock() {
		defer repo.modules.Unlock()
	}
	if !repo.Has(module, key) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return repo.Module(module).Get(key)
}

// Unmarshal unmarshal
func (repo *Repository) Unmarshal(module string, dest any) error {
	viper, ok := repo.modules.Get(module)
	if !ok {
		return NewModuleNotFoundException(module)
	}
	return viper.Unmarshal(dest)
}

// UnmarshalKey unmarshal key
func (repo *Repository) UnmarshalKey(module string, key string, dest any) error {
	viper, ok := repo.modules.Get(module)
	if !ok {
		return NewModuleNotFoundException(module)
	}
	return viper.UnmarshalKey(key, dest)
}
