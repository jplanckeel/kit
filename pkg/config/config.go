package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Tools map[string][]string

type Config map[string]Tools

func NewSource(sourceFile string) (Config, error) {
	var cfg Config
	source, err := os.ReadFile(filepath.Clean(sourceFile))
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to unmarshal %q: %w", sourceFile, err)
	}
	return cfg, nil

}

/*
func chartList(input []byte) ([]chartSearch, error) {
	var y []chartSearch
	err := yaml.Unmarshal(input, &y)
	if err != nil {
		return y, fmt.Errorf("failed to unmarshal %q: %w", input, err)
	}
	return y, nil

}*/
