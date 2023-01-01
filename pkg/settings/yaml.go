package settings

import (
	"github.com/dertseha/goconsider/pkg/consider"
	"gopkg.in/yaml.v3"
)

// FromYaml parses the provided raw YAML data into a settings instance.
func FromYaml(data []byte) (consider.Settings, error) {
	var settings consider.Settings
	err := yaml.Unmarshal(data, &settings)
	return settings, err
}
