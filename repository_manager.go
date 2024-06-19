package config

import (
	"strings"

	"github.com/gopi-frame/contract/config"
	"github.com/gopi-frame/support/maps"
	"github.com/mitchellh/mapstructure"
)

func NewRepositoryManager() *RepositoryManager {
	return &RepositoryManager{
		repos: maps.NewMap[string, config.Repository](),
	}
}

type RepositoryManager struct {
	repos *maps.Map[string, config.Repository]
}

func (rm *RepositoryManager) Repo(name string) config.Repository {
	return rm.repos.GetOr(name, NewRepository())
}

func (rm *RepositoryManager) Repos() map[string]config.Repository {
	return rm.repos.ToMap()
}

func (rm *RepositoryManager) AddRepo(name string, repo config.Repository) {
	rm.repos.Set(name, repo)
}

func (rm *RepositoryManager) All() map[string]any {
	values := make(map[string]any)
	rm.repos.Each(func(key string, value config.Repository) bool {
		values[key] = value.All()
		return true
	})
	return values
}

func (rm *RepositoryManager) Has(path string) bool {
	parts := strings.SplitN(path, ".", 2)
	if len(parts) == 1 {
		return rm.repos.ContainsKey(path)
	}
	repo, ok := rm.repos.Get(parts[0])
	if !ok {
		return false
	}
	return repo.Has(parts[1])
}

func (rm *RepositoryManager) Get(path string, defaultValue ...any) any {
	parts := strings.SplitN(path, ".", 2)
	if repo, ok := rm.repos.Get(parts[0]); ok {
		if len(parts) == 1 {
			return repo
		} else {
			return repo.Get(parts[1], defaultValue...)
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (rm *RepositoryManager) Set(path string, value any) {
	parts := strings.SplitN(path, ".", 2)
	if len(parts) == 1 {
		rm.AddRepo(path, value.(config.Repository))
	} else if repo, ok := rm.repos.Get(path); ok {
		repo.Set(parts[1], value)
	} else {
		repo := NewRepository()
		repo.Viper.Set(parts[1], value)
		rm.AddRepo(parts[0], repo)
	}
}

func (rm *RepositoryManager) Unmarshal(dest any) error {
	values := rm.All()
	return mapstructure.Decode(values, dest)
}

func (rm *RepositoryManager) UnmarshalKey(path string, dest any) error {
	parts := strings.SplitN(path, ".", 2)
	if repo, ok := rm.repos.Get(parts[0]); ok {
		if len(parts) == 1 {
			return repo.Unmarshal(dest)
		}
		return repo.UnmarshalKey(parts[1], dest)
	}
	return nil
}
