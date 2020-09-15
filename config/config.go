package config

import (
	"encoding/json"
	"fmt"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// PostgresConfig implements PostgreSQL configuration
type PostgresConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
}

// ConnectionString returns PostgreSQL connection string
func (c *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.Database,
	)
}

// Validate validates PostgreSQL configuration
func (c *PostgresConfig) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Host, validation.Required, is.Host),
		validation.Field(&c.Port, validation.Required, is.Port),
		validation.Field(&c.Username, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.Database, validation.Required),
	)
}

// ServerConfig implements API server configuration
type ServerConfig struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}

// Address returns API server address
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// Validate validates API server configuration
func (c *ServerConfig) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Host, validation.Required, is.Host),
		validation.Field(&c.Port, validation.Required, is.Port),
	)
}

// Config implements API configuration
type Config struct {
	LoggerConfig   zap.Config     `json:"logger" yaml:"logger"`
	DatabaseConfig PostgresConfig `json:"database" yaml:"database"`
	ServerConfig   ServerConfig   `json:"server" yaml:"server"`
}

// Validate validates API configuration
func (c *Config) Validate() error {
	if _, err := c.LoggerConfig.Build(); err != nil {
		return err
	}
	if err := c.DatabaseConfig.Validate(); err != nil {
		return err
	}
	if err := c.ServerConfig.Validate(); err != nil {
		return err
	}
	return nil
}

// LoadFromFile loads API configuration from a file located in the path
func (c *Config) LoadFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	switch ext := filepath.Ext(path); ext {
	case ".json":
		if err := json.Unmarshal(data, c); err != nil {
			return err
		}
	case ".yaml":
	case ".yml":
		if err := yaml.Unmarshal(data, c); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported configuration file format: %s", ext)
	}

	return c.Validate()
}
