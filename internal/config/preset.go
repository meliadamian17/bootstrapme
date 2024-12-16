package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type FileSpec struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
}

// Preset represents a configuration preset loaded from YAML.
type Preset struct {
	Name            string            `yaml:"name"`
	Description     string            `yaml:"description"`
	Language        string            `yaml:"language"`
	Framework       string            `yaml:"framework"`
	Variables       map[string]string `yaml:"variables"`
	PostInstallCmds []string          `yaml:"post_install_commands"`
	Files           []FileSpec        `yaml:"files"`
}

// LoadPresetFromFile loads a single preset from a given YAML file.
func LoadPresetFromFile(path string) (Preset, error) {
	var p Preset
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return p, err
	}
	err = yaml.Unmarshal(data, &p)
	return p, err
}
