package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	data   map[string]interface{}
	file   string
	plugin interface{}
}

func New(file string, plugin interface{}) *Config {
	return &Config{
		data:   make(map[string]interface{}),
		file:   file,
		plugin: plugin,
	}
}

func (c *Config) Load() error {
	// Si el archivo no existe, simplemente no cargamos nada.
	// Los defaults se aplican vía SetDefaults() desde el plugin.
	// El archivo se crea solo cuando el plugin llama explícitamente a Save().
	if _, err := os.Stat(c.file); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(c.file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &c.data)
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c.data)
	if err != nil {
		return err
	}

	dir := filepath.Dir(c.file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(c.file, data, 0644)
}

func (c *Config) Defaults() map[string]interface{} {
	return make(map[string]interface{})
}

func (c *Config) Get(key string) interface{} {
	return c.data[key]
}

func (c *Config) GetString(key string) string {
	if v, ok := c.data[key].(string); ok {
		return v
	}
	return ""
}

func (c *Config) GetInt(key string) int {
	if v, ok := c.data[key].(int); ok {
		return v
	}
	return 0
}

func (c *Config) GetBool(key string) bool {
	if v, ok := c.data[key].(bool); ok {
		return v
	}
	return false
}

func (c *Config) GetFloat(key string) float64 {
	if v, ok := c.data[key].(float64); ok {
		return v
	}
	return 0.0
}

func (c *Config) GetStringSlice(key string) []string {
	if v, ok := c.data[key].([]interface{}); ok {
		result := make([]string, len(v))
		for i, e := range v {
			if s, ok := e.(string); ok {
				result[i] = s
			}
		}
		return result
	}
	return nil
}

func (c *Config) GetMap(key string) map[string]interface{} {
	if v, ok := c.data[key].(map[string]interface{}); ok {
		return v
	}
	return make(map[string]interface{})
}

func (c *Config) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *Config) SetDefaults(defaults map[string]interface{}) {
	for k, v := range defaults {
		if _, exists := c.data[k]; !exists {
			c.data[k] = v
		}
	}
}

func (c *Config) All() map[string]interface{} {
	return c.data
}
