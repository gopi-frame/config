package config

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/knadh/koanf/v2"
)

var DefaultDelim = "."

var k *koanf.Koanf

var once sync.Once

func instance() *koanf.Koanf {
	once.Do(func() {
		if k == nil {
			k = koanf.New(DefaultDelim)
		}
	})
	return k
}

// New returns a new konaf instance
func New() *koanf.Koanf {
	return koanf.New(DefaultDelim)
}

// Use replaces the default instance
func Use(instance *koanf.Koanf) {
	k = instance
}

// Copy copies a [*koanf.Koanf] instance from the default one
func Copy() *koanf.Koanf {
	return instance().Copy()
}

// Cut creates [*koanf.Koanf] instance from given path
func Cut(path string) *koanf.Koanf {
	return instance().Cut(path)
}

// Load loads config
func Load(provider koanf.Provider, parser koanf.Parser, opts ...koanf.Option) error {
	return instance().Load(provider, parser, opts...)
}

// All returns all configs
func All() map[string]any {
	return instance().All()
}

// Exists returns whether the given path exists
func Exists(path string) bool {
	return instance().Exists(path)
}

// Get returns the value of a given path
func Get[T any](path string) T {
	var value T
	var convert = func(value any) T { return value.(T) }
	switch any(value).(type) {
	case int64:
		return convert(instance().Int64(path))
	case []int64:
		return convert(instance().Int64s(path))
	case map[string]int64:
		return convert(instance().Int64Map(path))
	case int:
		return convert(instance().Int(path))
	case []int:
		return convert(instance().Ints(path))
	case map[string]int:
		return convert(instance().IntMap(path))
	case float64:
		return convert(instance().Float64(path))
	case []float64:
		return convert(instance().Float64s(path))
	case map[string]float64:
		return convert(instance().Float64Map(path))
	case time.Duration:
		return convert(instance().Duration(path))
	case string:
		return convert(instance().String(path))
	case []string:
		return convert(instance().Strings(path))
	case map[string]string:
		return convert(instance().StringMap(path))
	case map[string][]string:
		return convert(instance().StringsMap(path))
	case []byte:
		return convert(instance().Bytes(path))
	case bool:
		return convert(instance().Bool(path))
	case []bool:
		return convert(instance().Bools(path))
	case map[string]bool:
		return convert(instance().BoolMap(path))
	default:
		if v, ok := instance().Get(path).(T); ok {
			return v
		}
		return value
	}
}

// Must returns the value of a given path or panics if the given path doesn't exist or it's value is empty
func Must[T any](path string) T {
	var value T
	var convert = func(value any) T { return value.(T) }
	switch any(value).(type) {
	case int64:
		return convert(instance().MustInt64(path))
	case []int64:
		return convert(instance().MustInt64s(path))
	case map[string]int64:
		return convert(instance().MustInt64Map(path))
	case int:
		return convert(instance().MustInt(path))
	case []int:
		return convert(instance().MustInts(path))
	case map[string]int:
		return convert(instance().MustIntMap(path))
	case float64:
		return convert(instance().MustFloat64(path))
	case []float64:
		return convert(instance().MustFloat64s(path))
	case map[string]float64:
		return convert(instance().MustFloat64Map(path))
	case time.Duration:
		return convert(instance().MustDuration(path))
	case string:
		return convert(instance().MustString(path))
	case []string:
		return convert(instance().Strings(path))
	case map[string]string:
		return convert(instance().MustStringMap(path))
	case map[string][]string:
		return convert(instance().MustStringsMap(path))
	case []byte:
		return convert(instance().MustBytes(path))
	case bool:
		return convert(instance().Bool(path))
	case []bool:
		return convert(instance().MustBools(path))
	case map[string]bool:
		return convert(instance().MustBoolMap(path))
	default:
		if v, ok := instance().Get(path).(T); ok {
			if reflect.ValueOf(v).IsZero() {
				panic(fmt.Sprintf("invalid value %s=%v", path, v))
			}
			return v
		}
		panic(fmt.Sprintf("invalid value %s=%v", path, value))
	}
}

// Time returns the value of a given path and attempts to parse it into [time.Time]
func Time(path string, layout string) time.Time {
	return instance().Time(path, layout)
}

// MustTime returns the value of a given path and attempts to parse it into [time.Time], it panics when the time is zero
func MustTime(path string, layout string) time.Time {
	return instance().MustTime(path, layout)
}

// Unmarshal attempts to unmarshal the value of a given path
func Unmarshal(path string, output any) error {
	return instance().Unmarshal(path, output)
}

// UnmarshalWithConf attempts to unmarshal the value of a given path with custom config
func UnmarshalWithConf(path string, output any, config koanf.UnmarshalConf) error {
	return instance().UnmarshalWithConf(path, output, config)
}

// Marshal marshals the config map into bytes
func Marshal(p koanf.Parser) ([]byte, error) {
	return instance().Marshal(p)
}
