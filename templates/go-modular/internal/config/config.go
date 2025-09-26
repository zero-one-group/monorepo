package config

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var globalConfig *Config

// Load configuration from multiple sources with priority:
// 1. Environment variables (highest priority)
// 2. .env file (if exists)
// 3. Default values (lowest priority)
func Load(configEnvPath string) (*Config, error) {
	loadDotEnv(configEnvPath) // Load .env file first (if exists)

	// Only use environment variables and defaults. No config file support.
	viper.AutomaticEnv()

	// Set all defaults from DefaultConfig using lowercase struct field names
	setViperDefaults("", DefaultConfig())

	// Bind envs for all config fields using `env` tag values
	bindEnvRecursive("", reflect.TypeOf(Config{}))

	// Unmarshal into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate critical configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	globalConfig = &config
	return &config, nil
}

// Get returns the global configuration instance
func Get() *Config {
	if globalConfig == nil {
		panic("configuration not loaded. Call config.Load() first")
	}
	return globalConfig
}

// Loads environment variables from .env file or a provided env file path
func loadDotEnv(path string) {
	if path != "" {
		_ = godotenv.Load(path)
		return
	}
	// Try to load .env file from working directory
	if err := godotenv.Load(); err != nil {
		// Try to load from config directory
		_ = godotenv.Load("config/.env")
	}
}

// setViperDefaults sets viper defaults recursively from struct values.
// Keys are built using lowercase field names joined by '.' (e.g. "app.mode").
func setViperDefaults(prefix string, val any) {
	rv := reflect.ValueOf(val)
	rt := reflect.TypeOf(val)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return
		}
		rv = rv.Elem()
		rt = rt.Elem()
	}
	if rv.Kind() == reflect.Struct {
		for i := 0; i < rv.NumField(); i++ {
			fieldVal := rv.Field(i)
			fieldType := rt.Field(i)
			keyPart := strings.ToLower(fieldType.Name)
			var key string
			if prefix == "" {
				key = keyPart
			} else {
				key = prefix + "." + keyPart
			}
			// Recurse into structs
			if fieldVal.Kind() == reflect.Struct || (fieldVal.Kind() == reflect.Pointer && fieldVal.Elem().Kind() == reflect.Struct) {
				setViperDefaults(key, fieldVal.Interface())
			} else {
				viper.SetDefault(key, fieldVal.Interface())
			}
		}
	} else {
		// For non-struct leaf values, set default for the current prefix
		if prefix != "" {
			viper.SetDefault(prefix, val)
		}
	}
}

// bindEnvRecursive binds environment variables for all fields with `env` tag recursively.
// Key names use lowercase field names joined by '.' to match defaults and nested maps.
func bindEnvRecursive(prefix string, t reflect.Type) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		keyPart := strings.ToLower(field.Name)
		var key string
		if prefix == "" {
			key = keyPart
		} else {
			key = prefix + "." + keyPart
		}

		// env tag may be empty or contain ",squash"
		envTag := field.Tag.Get("env")
		// If envTag contains a real env var name (not just squash), bind it
		if envTag != "" && envTag != ",squash" {
			// env tag can include options after comma, take first part before comma
			envName := strings.SplitN(envTag, ",", 2)[0]
			if envName != "" {
				_ = viper.BindEnv(key, envName)
			}
		}

		// Recurse into nested structs
		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Struct {
			bindEnvRecursive(key, ft)
		}
	}
}

// GenerateExampleEnvFile writes a .env.example file with all config keys and their default values,
// grouped and sorted by top-level config section order.
func GenerateExampleEnvFile(path string) error {
	cfg := DefaultConfig()
	envs := make(map[string]map[string]string)
	sectionOrder := []string{}

	cfgVal := reflect.ValueOf(cfg)
	cfgType := reflect.TypeOf(cfg)

	// Loop all top-level fields in Config struct
	for i := 0; i < cfgVal.NumField(); i++ {
		section := cfgType.Field(i).Name
		sectionOrder = append(sectionOrder, section)
		envs[section] = make(map[string]string)

		// Extract env tags for each section
		var extractEnv func(val any)
		extractEnv = func(val any) {
			rv := reflect.ValueOf(val)
			rt := reflect.TypeOf(val)
			if rv.Kind() == reflect.Pointer {
				if rv.IsNil() {
					return
				}
				rv = rv.Elem()
				rt = rt.Elem()
			}
			if rv.Kind() == reflect.Struct {
				for j := 0; j < rv.NumField(); j++ {
					field := rv.Field(j)
					fieldType := rt.Field(j)
					envTag := fieldType.Tag.Get("env")
					if envTag != "" && envTag != "-" && !strings.Contains(envTag, "squash") {
						envName := strings.SplitN(envTag, ",", 2)[0]
						if envName != "" {
							if _, exists := envs[section][envName]; !exists {
								envs[section][envName] = fmt.Sprintf("%v", field.Interface())
							}
						}
					}
					if field.Kind() == reflect.Struct || (field.Kind() == reflect.Pointer && field.Elem().Kind() == reflect.Struct) {
						extractEnv(field.Interface())
					}
				}
			}
		}
		extractEnv(cfgVal.Field(i).Interface())
	}

	var buf bytes.Buffer
	buf.WriteString("# Example Application Environment Configuration\n")
	buf.WriteString("# Copy this file to .env and update the values\n\n")

	for _, section := range sectionOrder {
		keys := make([]string, 0, len(envs[section]))
		for k := range envs[section] {
			keys = append(keys, k)
		}
		if len(keys) == 0 {
			continue
		}
		sort.Strings(keys)
		buf.WriteString(fmt.Sprintf("# %s\n", section))
		for _, k := range keys {
			buf.WriteString(fmt.Sprintf("%s=%v\n", k, envs[section][k]))
		}
		buf.WriteString("\n")
	}

	return os.WriteFile(path, buf.Bytes(), 0600)
}
