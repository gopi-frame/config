package config

import (
	"github.com/spf13/viper"
)

func NewRepository() *Repository {
	repo := new(Repository)
	repo.Viper = viper.New()
	return repo
}

// Repository repository
type Repository struct {
	*viper.Viper
}

func (repo *Repository) All() map[string]any {
	return repo.Viper.AllSettings()
}

// Has has
func (repo *Repository) Has(path string) bool {
	return repo.IsSet(path)
}

// Get get
func (repo *Repository) Get(path string, defaultValue ...any) any {
	if repo.Has(path) {
		return repo.Viper.Get(path)
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (repo *Repository) Set(path string, value any) {
	repo.Viper.Set(path, value)
}

// Unmarshal unmarshal
func (repo *Repository) Unmarshal(dest any) error {
	return repo.Viper.Unmarshal(dest)
}

func (repo *Repository) UnmarshalKey(key string, dest any) error {
	return repo.Viper.UnmarshalKey(key, dest)
}
