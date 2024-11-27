package config

import (
	"strings"
	"testing"

	"github.com/knadh/koanf/maps"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Load(t *testing.T) {
	repo := NewRepository()
	err := repo.Load(ProviderFunc(func() ([]byte, error) {
		return []byte(`user=root,password=root,database=gopi,host=localhost,port=3306,params.charset=utf8mb4`), nil
	}), ParserFunc(func(data []byte) (map[string]any, error) {
		str := string(data)
		pairs := strings.Split(str, ",")
		result := make(map[string]any)
		for _, pair := range pairs {
			keyValue := strings.Split(pair, "=")
			if len(keyValue) == 2 {
				key := keyValue[0]
				value := keyValue[1]
				result[key] = value
			}
		}
		return maps.Unflatten(result, "."), nil
	}))
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, "root", repo.String("user"))
	assert.Equal(t, "root", repo.String("password"))
	assert.Equal(t, "localhost", repo.String("host"))
	assert.Equal(t, "3306", repo.String("port"))
	assert.Equal(t, "utf8mb4", repo.String("params.charset"))
	assert.Equal(t, map[string]any{"charset": "utf8mb4"}, repo.AnyMap("params"))
}

func TestRepository_Set(t *testing.T) {
	repo := NewRepository()
	err := repo.Set("user", "root")
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, "root", repo.String("user"))
}

func TestRepository_Delete(t *testing.T) {
	repo := NewRepository()
	err := repo.Set("user", "root")
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, "root", repo.String("user"))
	err = repo.Delete("user")
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, "", repo.String("user"))
}
