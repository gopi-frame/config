package yaml

import "gopkg.in/yaml.v3"

type YamlParser struct{}

func NewYamlParser() *YamlParser {
	return &YamlParser{}
}

func (p *YamlParser) Unmarshal(data []byte) (map[string]any, error) {
	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return config, nil
}
