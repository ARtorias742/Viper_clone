package myviper

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// Viper is the main configuration struct
type Viper struct {
	configName string
	configType string
	configPath []string
	keyValue   map[string]interface{}
	fs         afero.Fs
}

// New creates a new Viper instance
func New() *Viper {
	return &Viper{
		keyValue: make(map[string]interface{}),
		fs:       afero.NewOsFs(),
	}
}

// SetConfigName sets the name of the config file (without extension)
func (v *Viper) SetConfigName(name string) {
	v.configName = name
}

// SetConfigType sets the type of the config file (e.g., "yaml")
func (v *Viper) SetConfigType(typ string) {
	v.configType = typ
}

// AddConfigPath adds a path to search for the config file
func (v *Viper) AddConfigPath(path string) {
	v.configPath = append(v.configPath, path)
}

// ReadInConfig reads the configuration file into memory
func (v *Viper) ReadInConfig() error {
	file, err := v.findConfigFile()
	if err != nil {
		return err
	}

	data, err := afero.ReadFile(v.fs, file)
	if err != nil {
		return err
	}

	switch strings.ToLower(v.configType) {
	case "yaml", "yml":
		return v.unmarshalYAML(data)
	default:
		return fmt.Errorf("unsupported config type: %s", v.configType)
	}
}

// findConfigFile searches for the config file in the specified paths
func (v *Viper) findConfigFile() (string, error) {
	for _, path := range v.configPath {
		fullPath := filepath.Join(path, v.configName+"."+v.configType)
		exists, err := afero.Exists(v.fs, fullPath)
		if err == nil && exists {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("config file not found: %s", v.configName)
}

// unmarshalYAML unmarshals YAML data into the keyValue map
func (v *Viper) unmarshalYAML(data []byte) error {
	return yaml.Unmarshal(data, &v.keyValue)
}

// Get retrieves a value by key
func (v *Viper) Get(key string) interface{} {
	return v.keyValue[key]
}

// GetString retrieves a value as a string
func (v *Viper) GetString(key string) string {
	if val, ok := v.keyValue[key].(string); ok {
		return val
	}
	return ""
}

// Set sets a key-value pair
func (v *Viper) Set(key string, value interface{}) {
	v.keyValue[key] = value
}

// WriteConfig writes the current configuration to the file
func (v *Viper) WriteConfig() error {
	file := filepath.Join(v.configPath[0], v.configName+"."+v.configType)
	var buf bytes.Buffer
	if err := yaml.NewEncoder(&buf).Encode(v.keyValue); err != nil {
		return err
	}
	return afero.WriteFile(v.fs, file, buf.Bytes(), 0644)
}
