package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type EnvType string

const (
	Development EnvType = "development"
	Production  EnvType = "production"
	Test        EnvType = "test"
)

func (env EnvType) isValid() bool {
	switch env {
	case Development, Production, Test:
		return true
	default:
		return false
	}
}

type Service struct {
	Port      int     `env:"PORT,default=8080"`
	Env       EnvType `env:"ENV,default=development"`
	SecretKey string  `env:"SECRET_KEY,default=very_secret"`
}

type Database struct {
	DbHost string `env:"DB_HOST,default=localhost"`
	DbUser string `env:"DB_USER,default=root"`
	DbPass string `env:"DB_PASS"`
	DbName string `env:"DB_NAME,default=test_db"`
	DbPort string `env:"DB_PORT,default=5432"`
}

type Midtrans struct {
	MidtransServiceKey string `env:"MIDTRANS_SERVER_KEY"`
	MidtransClientKey  string `env:"MIDTRANS_CLIENT_KEY"`
}

type Config struct {
	Service  Service
	Database Database
	Midtrans Midtrans
}

// LoadEnv loads environment variables and maps them to the Config struct.
func LoadEnv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var config Config

	err = loadFromEnv(&config)
	if err != nil {
		return nil, err
	}

	if !config.Service.Env.isValid() {
		return nil, fmt.Errorf("invalid env value: %s", config.Service.Env)
	}

	return &config, nil
}

// loadFromEnv populates the fields in a struct based on environment variables using reflection.
func loadFromEnv(config interface{}) error {
	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			// Recursive call for nested structs
			err := loadFromEnv(field.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}

		tag := fieldType.Tag.Get("env")
		if tag == "" {
			continue
		}

		tagParts := parseTag(tag)
		envVar := tagParts["env"]
		defaultValue := tagParts["default"]

		envValue := os.Getenv(envVar)
		if envValue == "" {
			envValue = defaultValue
		}

		if envValue == "" && field.Kind() != reflect.String {
			return fmt.Errorf("environment variable %s is required but not set", envVar)
		}

		// Set the field value
		err := setFieldValue(field, envValue)
		if err != nil {
			return err
		}
	}

	return nil
}

// setFieldValue sets the value of a struct field based on the provided string value.
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid int value: %v", err)
		}
		field.SetInt(int64(intVal))
	case reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int64 value: %v", err)
		}
		field.SetInt(intVal)
	case reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float64 value: %v", err)
		}
		field.SetFloat(floatVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

// parseTag parses a struct tag into a map containing the tag's key-value pairs.
func parseTag(tag string) map[string]string {
	result := make(map[string]string)
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		} else {
			result["env"] = kv[0]
		}
	}
	return result
}
