package config

import (
	"github.com/gopi-frame/contract/config"
	"github.com/gopi-frame/support/maps"
	"github.com/spf13/viper"
)

var _ config.Repository = (*Repository)(nil)

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
	repo.modules.TryLock()
	defer repo.modules.Unlock()
	return repo.modules.GetOr(module, nil)
}

// Set set
func (repo *Repository) Set(module string, key string, value any) {
	repo.modules.TryLock()
	defer repo.modules.Unlock()
	repo.Module(module).Set(key, value)
}

// Has has
func (repo *Repository) Has(module string, key string) bool {
	repo.modules.TryLock()
	defer repo.modules.Unlock()
	if module := repo.Module(module); module != nil {
		return module.IsSet(key)
	}
	return false
}

// Get get
func (repo *Repository) Get(module string, key string, defaultValue ...any) any {
	if !repo.Has(module, key) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return repo.Module(module).Get(key)
}
