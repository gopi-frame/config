package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/contract/repository"
	"github.com/knadh/koanf/maps"
)

var _ repository.Repository = (*Repository)(nil)

type Repository struct {
	mu     sync.RWMutex
	values map[string]any
}

func NewRepository() *Repository {
	return &Repository{
		values: make(map[string]any),
	}
}

func (r *Repository) get(path string) (any, error) {
	keys := strings.Split(path, ".")
	var value any = r.values
	for _, key := range keys {
		if m, ok := value.(map[string]any); ok {
			if v, found := m[key]; found {
				value = v
				continue
			} else {
				return nil, errors.New("key not found")
			}
		} else if s, ok := value.([]any); ok {
			if index, err := strconv.Atoi(key); err == nil {
				if index >= 0 && index < len(s) {
					value = s[index]
					continue
				} else {
					return nil, errors.New("index out of range")
				}
			} else {
				return nil, errors.New("invalid index")
			}
		} else {
			return nil, errors.New("invalid type")
		}
	}
	return value, nil
}

func (r *Repository) set(data map[string]any, keys []string, index int, value any) error {
	if index >= len(keys) {
		return errors.New("path is too short to set value")
	}
	key := keys[index]
	if index < len(keys)-1 {
		if _, exists := data[key]; !exists {
			nextKey := keys[index+1]
			if _, err := strconv.Atoi(nextKey); err == nil {
				data[key] = []any{}
			} else {
				data[key] = make(map[string]any)
			}
		}
		if m, ok := data[key].(map[string]any); ok {
			return r.set(m, keys, index+1, value)
		} else if s, ok := data[key].([]any); ok {
			if nextIndex, err := strconv.Atoi(keys[index+1]); err == nil {
				for len(s) <= nextIndex {
					s = append(s, []any{})
				}
				if _, ok := s[nextIndex].(map[string]any); !ok {
					s[nextIndex] = make(map[string]any)
				}
				return r.set(s[nextIndex].(map[string]any), keys, index+1, value)
			} else {
				return fmt.Errorf("cannot access slice with non-integer key %q", keys[index+1])
			}
		} else {
			return fmt.Errorf("unexpected value type at key %q", key)
		}
	}
	data[key] = value
	return nil
}

func (r *Repository) delete(data map[string]any, keys []string, index int) error {
	if index >= len(keys) {
		return errors.New("path is too short to delete key")
	}
	key := keys[index]
	if index < len(keys)-1 {
		if nestedMap, ok := data[key].(map[string]any); ok {
			return r.delete(nestedMap, keys, index+1)
		} else {
			return fmt.Errorf("cannot delete key %q: not a map", key)
		}
	}
	if _, exists := data[key]; exists {
		delete(data, key)
		return nil
	}
	return fmt.Errorf("key %q not found in map", key)
}

func (r *Repository) Load(provider repository.Provider, parser repository.Parser) error {
	content, err := provider.Read()
	if err != nil {
		return err
	}
	data, err := parser.Unmarshal(content)
	if err != nil {
		return err
	}
	for key, value := range data {
		r.values[key] = value
	}
	return nil
}

func (r *Repository) LoadAt(path string, provider repository.Provider, parser repository.Parser) error {
	content, err := provider.Read()
	if err != nil {
		return err
	}
	data, err := parser.Unmarshal(content)
	if err != nil {
		return err
	}
	return r.Set(path, data)
}

func (r *Repository) Has(path string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, err := r.get(path); err == nil {
		return true
	}
	return false
}

func (r *Repository) Get(path string) (any, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.get(path)
}

func (r *Repository) Set(path string, value any) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	keys := strings.Split(path, ".")
	return r.set(r.values, keys, 0, value)
}

func (r *Repository) Delete(path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	keys := strings.Split(path, ".")
	return r.delete(r.values, keys, 0)
}

func (r *Repository) Keys() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var keys []string
	for key := range r.values {
		keys = append(keys, key)
	}
	return keys
}

func (r *Repository) Paths() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var paths []string
	_, keyMaps := maps.Flatten(r.values, nil, ".")
	for _, path := range keyMaps {
		paths = append(paths, strings.Join(path, "."))
	}
	return paths
}

func (r *Repository) ToMap() map[string]any {
	return r.values
}

func (r *Repository) Cut(path string) repository.Repository {
	r.mu.RLock()
	defer r.mu.RUnlock()
	repo := NewRepository()
	if value, err := r.get(path); err == nil {
		if value, ok := value.(map[string]any); ok {
			repo.values = value
		}
	}
	return repo
}

func (r *Repository) Merge(repo repository.Repository) (repository.Repository, error) {
	for key, value := range repo.ToMap() {
		if err := r.Set(key, value); err != nil {
			return r, err
		}
	}
	return r, nil
}

func (r *Repository) MergeAt(path string, repo repository.Repository) (repository.Repository, error) {
	if value, err := r.get(path); err == nil {
		for key, value := range repo.ToMap() {
			if err := r.Set(path+"."+key, value); err != nil {
				return r, err
			}
		}
	} else {
		if err := r.Set(path, value); err != nil {
			return r, err
		}
	}
	return r, nil
}

func (r *Repository) Unmarshal(path string, v any, conf *mapstructure.DecoderConfig) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if conf.Result == nil {
		conf.Result = v
	}
	decoder, err := mapstructure.NewDecoder(conf)
	if err != nil {
		return err
	}
	if values, err := r.get(path); err == nil {
		if err := decoder.Decode(values); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Int64(path string, defaultValue ...int64) int64 {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(int64); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (r *Repository) Int64s(path string, defaultValue ...[]int64) []int64 {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.([]int64); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) Int64Map(path string, defaultValue ...map[string]int64) map[string]int64 {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(map[string]int64); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) Int(path string, defaultValue ...int) int {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(int); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (r *Repository) Ints(path string, defaultValue ...[]int) []int {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.([]int); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) IntMap(path string, defaultValue ...map[string]int) map[string]int {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(map[string]int); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) Float64(path string, defaultValue ...float64) float64 {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(float64); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (r *Repository) Float64s(path string, defaultValue ...[]float64) []float64 {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.([]float64); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) Float64Map(path string, defaultValue ...map[string]float64) map[string]float64 {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(map[string]float64); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) Duration(path string, defaultValue ...time.Duration) time.Duration {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(time.Duration); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (r *Repository) Time(path string, layout string, defaultValue ...time.Time) time.Time {
	if value, err := r.Get(path); err == nil {
		if v, ok := value.(time.Time); ok {
			return v
		} else if v, ok := value.(string); ok {
			if t, err := time.Parse(layout, v); err == nil {
				return t
			}
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return time.Time{}
}

func (r *Repository) String(path string, defaultValue ...string) string {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(string); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (r *Repository) Strings(path string, defaultValue ...[]string) []string {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.([]string); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) StringMap(path string, defaultValue ...map[string]string) map[string]string {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(map[string]string); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) StringsMap(path string, defaultValue ...map[string][]string) map[string][]string {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(map[string][]string); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) Bool(path string, defaultValue ...bool) bool {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(bool); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

func (r *Repository) Bools(path string, defaultValue ...[]bool) []bool {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.([]bool); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) BoolMap(path string, defaultValue ...map[string]bool) map[string]bool {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(map[string]bool); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) Bytes(path string, defaultValue ...[]byte) []byte {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.([]byte); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) Any(path string, defaultValue ...any) any {
	if value, err := r.Get(path); err == nil {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) Anys(path string, defaultValue ...[]any) []any {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.([]any); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (r *Repository) AnyMap(path string, defaultValue ...map[string]any) map[string]any {
	if value, err := r.Get(path); err == nil {
		if value, ok := value.(map[string]any); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}
