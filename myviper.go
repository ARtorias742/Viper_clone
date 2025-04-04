package myviper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// Viper is the main configuration struct
type Viper struct {
	configName string
	configType string
	configPath []string
	keyValue   map[string]interface{}
	fs         afero.Fs
	flags      *pflag.FlagSet
	envEnabled bool   // Toggle for environment variable binding
	envPrefix  string // Optional prefix for env vars (e.g., "APP_")
}

// New creates a new Viper instance
func New() *Viper {
	v := &Viper{
		keyValue:   make(map[string]interface{}),
		fs:         afero.NewOsFs(),
		flags:      pflag.NewFlagSet("myviper", pflag.ExitOnError),
		envEnabled: false, // Env binding off by default
		envPrefix:  "",    // No prefix by default
	}
	return v
}

// SetConfigName sets the name of the config file (without extension)
func (v *Viper) SetConfigName(name string) {
	v.configName = name
}

// SetConfigType sets the type of the config file (e.g., "yaml", "json")
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
	case "json":
		return v.unmarshalJSON(data)
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

// unmarshalJSON unmarshals JSON data into the keyValue map
func (v *Viper) unmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.keyValue)
}

// Get retrieves a value by key
func (v *Viper) Get(key string) interface{} {
	// Check flags first, then keyValue, then environment
	if v.flags != nil {
		if flagVal, err := v.flags.GetString(key); err == nil && flagVal != "" {
			return flagVal
		}
	}
	if val, ok := v.keyValue[key]; ok {
		return val
	}
	if v.envEnabled {
		envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		if v.envPrefix != "" {
			envKey = v.envPrefix + "_" + envKey
		}
		if envVal := os.Getenv(envKey); envVal != "" {
			return envVal
		}
	}
	return nil
}

// GetString retrieves a value as a string
func (v *Viper) GetString(key string) string {
	if val := v.Get(key); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
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
	switch strings.ToLower(v.configType) {
	case "yaml", "yml":
		if err := yaml.NewEncoder(&buf).Encode(v.keyValue); err != nil {
			return err
		}
	case "json":
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "  ")
		if err := enc.Encode(v.keyValue); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported config type: %s", v.configType)
	}
	return afero.WriteFile(v.fs, file, buf.Bytes(), 0644)
}

// BindPFlags binds command-line flags to Viper
func (v *Viper) BindPFlags(flags *pflag.FlagSet) {
	v.flags = flags
}

// AutomaticEnv enables automatic binding of environment variables
func (v *Viper) AutomaticEnv() {
	v.envEnabled = true
}

// SetEnvPrefix sets a prefix for environment variables (e.g., "APP")
func (v *Viper) SetEnvPrefix(prefix string) {
	v.envPrefix = strings.ToUpper(prefix)
}
