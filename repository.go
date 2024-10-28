package config

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/contract/repository"
	"github.com/knadh/koanf/maps"
	"strings"
	"time"
)

var _ repository.Repository = (*Repository)(nil)

type Repository struct {
	values *kv.Map[string, any]
	keyMap map[string][]string
}

func NewRepository() *Repository {
	return &Repository{
		values: kv.NewMap[string, any](),
	}
}

func (r *Repository) findKeyStartsWith(path string) []string {
	keys := make([]string, 0)
	for key := range r.keyMap {
		if strings.HasPrefix(key, path+".") {
			keys = append(keys, key)
		}
	}
	return keys
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
	flatten, keyMap := maps.Flatten(data, nil, ".")
	for key, value := range flatten {
		r.values.Set(key, value)
	}
	for key, value := range keyMap {
		r.keyMap[key] = value
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
	data = maps.Unflatten(map[string]any{
		path: data,
	}, ".")
	flatten, keyMap := maps.Flatten(data, nil, ".")
	for key, value := range flatten {
		r.values.Set(key, value)
	}
	for key, value := range keyMap {
		r.keyMap[key] = value
	}
	return nil
}

func (r *Repository) Has(path string) bool {
	return r.values.Contains(path)
}

func (r *Repository) Get(path string) any {
	return r.values.GetOr(path, nil)
}

func (r *Repository) Set(path string, value any) {
	if r.Has(path) {
		r.values.Set(path, value)
		return
	}
	removeKeys := r.findKeyStartsWith(path)
	for _, key := range removeKeys {
		r.values.Remove(key)
		delete(r.keyMap, key)
	}
	data := maps.Unflatten(map[string]any{
		path: value,
	}, ".")
	flatten, mapKey := maps.Flatten(data, nil, ".")
	for key, value := range flatten {
		r.values.Set(key, value)
	}
	for key, value := range mapKey {
		r.keyMap[key] = value
	}
}

func (r *Repository) Delete(path string) {
	if r.Has(path) {
		r.values.Remove(path)
		delete(r.keyMap, path)
		return
	}
	removeKeys := r.findKeyStartsWith(path)
	for _, key := range removeKeys {
		r.values.Remove(key)
		delete(r.keyMap, key)
	}
}

func (r *Repository) Keys() []string {
	keys := make([]string, 0)
	for key := range r.keyMap {
		keys = append(keys, key)
	}
	return keys
}

func (r *Repository) All() map[string]any {
	return r.values.ToMap()
}

func (r *Repository) Cut(path string) repository.Repository {
	repo := NewRepository()
	if r.Has(path) {
		repo.Set(path, r.Get(path))
		return repo
	}
	keys := r.findKeyStartsWith(path)
	for _, key := range keys {
		repo.Set(key, r.Get(path))
		repo.keyMap[key] = r.keyMap[key]
	}
	return repo
}

func (r *Repository) Merge(repo repository.Repository) {
	for _, key := range repo.Keys() {
		r.Set(key, repo.Get(key))
	}
}

func (r *Repository) MergeAt(path string, repo repository.Repository) {
	for _, key := range repo.Keys() {
		r.Set(path+"."+key, repo.Get(key))
	}
}

func (r *Repository) Unmarshal(path string, v any, conf *mapstructure.DecoderConfig) error {
	if conf.Result == nil {
		conf.Result = v
	}
	decoder, err := mapstructure.NewDecoder(conf)
	if err != nil {
		return err
	}
	if r.Has(path) {
		return decoder.Decode(r.Get(path))
	}
	var values = map[string]any{}
	keys := r.findKeyStartsWith(path)
	for _, key := range keys {
		values[key] = r.Get(key)
	}
	return decoder.Decode(maps.Unflatten(values, "."))
}

func (r *Repository) Int64(path string, defaultValue ...int64) int64 {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	value, ok := r.Get(path).(int64)
	if !ok {
		return 0
	}
	return value
}

func (r *Repository) Int64s(path string, defaultValue ...[]int64) []int64 {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).([]int64)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) Int64Map(path string, defaultValue ...map[string]int64) map[string]int64 {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).(map[string]int64)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) Int(path string, defaultValue ...int) int {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	value, ok := r.Get(path).(int)
	if !ok {
		return 0
	}
	return value
}

func (r *Repository) Ints(path string, defaultValue ...[]int) []int {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []int{}
	}
	value, ok := r.Get(path).([]int)
	if !ok {
		return []int{}
	}
	return value
}

func (r *Repository) IntMap(path string, defaultValue ...map[string]int) map[string]int {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).(map[string]int)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) Float64(path string, defaultValue ...float64) float64 {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	value, ok := r.Get(path).(float64)
	if !ok {
		return 0
	}
	return value
}

func (r *Repository) Float64s(path string, defaultValue ...[]float64) []float64 {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).([]float64)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) Float64Map(path string, defaultValue ...map[string]float64) map[string]float64 {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).(map[string]float64)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) Duration(path string, defaultValue ...time.Duration) time.Duration {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	value, ok := r.Get(path).(time.Duration)
	if !ok {
		return 0
	}
	return value
}

func (r *Repository) Time(path string, layout string, defaultValue ...time.Time) time.Time {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return time.Time{}
	}
	value, ok := r.Get(path).(string)
	if !ok {
		return time.Time{}
	}
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (r *Repository) String(path string, defaultValue ...string) string {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	value, ok := r.Get(path).(string)
	if !ok {
		return ""
	}
	return value
}

func (r *Repository) Strings(path string, defaultValue ...[]string) []string {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).([]string)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) StringMap(path string, defaultValue ...map[string]string) map[string]string {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).(map[string]string)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) StringsMap(path string, defaultValue ...map[string][]string) map[string][]string {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).(map[string][]string)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) Bool(path string, defaultValue ...bool) bool {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	value, ok := r.Get(path).(bool)
	if !ok {
		return false
	}
	return value
}

func (r *Repository) Bools(path string, defaultValue ...[]bool) []bool {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).([]bool)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) BoolMap(path string, defaultValue ...map[string]bool) map[string]bool {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).(map[string]bool)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) Bytes(path string, defaultValue ...[]byte) []byte {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).([]byte)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) Any(path string, defaultValue ...any) any {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return r.Get(path)
}

func (r *Repository) Anys(path string, defaultValue ...[]any) []any {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).([]any)
	if !ok {
		return nil
	}
	return value
}

func (r *Repository) AnyMap(path string, defaultValue ...map[string]any) map[string]any {
	if !r.Has(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	value, ok := r.Get(path).(map[string]any)
	if !ok {
		return nil
	}
	return value
}
