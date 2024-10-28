package config

type ProviderFunc func() ([]byte, error)

func (f ProviderFunc) Read() ([]byte, error) {
	return f()
}
