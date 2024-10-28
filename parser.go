package config

type ParserFunc func(data []byte) (map[string]any, error)

func (f ParserFunc) Unmarshal(data []byte) (map[string]any, error) {
	return f(data)
}
