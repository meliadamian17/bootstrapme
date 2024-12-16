package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func LoadAllPresets() (map[string][]Preset, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".config", "bootstrapme")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return nil, nil
	}

	presetsByLang := make(map[string][]Preset)
	languages, err := ioutil.ReadDir(configDir)
	if err != nil {
		return nil, err
	}

	for _, langDir := range languages {
		if !langDir.IsDir() {
			continue
		}
		langPath := filepath.Join(configDir, langDir.Name())
		yamlFiles, err := ioutil.ReadDir(langPath)
		if err != nil {
			log.Printf("Error reading dir %s: %v", langPath, err)
			continue
		}

		for _, yf := range yamlFiles {
			if strings.HasSuffix(yf.Name(), ".yaml") {
				presetPath := filepath.Join(langPath, yf.Name())
				p, err := LoadPresetFromFile(presetPath)
				if err == nil && p.Language != "" && p.Framework != "" {
					presetsByLang[langDir.Name()] = append(presetsByLang[langDir.Name()], p)
				} else if err != nil {
					log.Printf("Error loading preset from %s: %v", presetPath, err)
				}
			}
		}
	}

	return presetsByLang, nil
}
