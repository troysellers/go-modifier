package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	SF       SFConfig
	Mockaroo MockarooConfig
}
type MockarooConfig struct {
	Key     string
	DataDir string
}
type SFConfig struct {
	Username    string
	Password    string
	Token       string
	LoginUrl    string
	ApiVersion  float32
	SfDebug     bool
	Queries     []string
	SfBatchSize int
}

// get the configuration from the environment variables.
func NewConfig() *Config {

	return &Config{
		SF: SFConfig{
			Username:    getEnv("SF_USER", ""),
			Password:    getEnv("SF_PASS", ""),
			Token:       getEnv("SF_TOKEN", ""),
			LoginUrl:    getEnv("SF_ENDPOINT", ""),
			ApiVersion:  getEnvFloat("SF_API_VERSION", 52.0),
			SfDebug:     getEnvBool("SF_DEBUG", false),
			Queries:     getEnvStringArray("QUERIES", ";"),
			SfBatchSize: getEnvInt("SF_BATCH_SIZE", 200),
		},
		Mockaroo: MockarooConfig{
			Key:     getEnv("MOCKAROO_KEY", ""),
			DataDir: getEnv("MOCKAROO_DATA_DIR", ""),
		},
	}
}

// get string environment variable
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// get boolean environment variable
func getEnvBool(key string, defaultVal bool) bool {
	envBool := getEnv(key, "false")
	b, err := strconv.ParseBool(envBool)
	if err != nil {
		return defaultVal
	}
	return b
}

// returns nil if key not found
// the separator defines the character to split on
func getEnvStringArray(key string, separator string) []string {
	s := getEnv(key, "")
	if s == key {
		return nil
	}
	return strings.Split(s, separator)
}

// returns an integer value for the key
func getEnvInt(key string, defaultVal int) int {
	s := getEnv(key, "")
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return i
}

func getEnvFloat(key string, defaultVal float32) float32 {
	s := getEnv(key, "")
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return defaultVal
	}
	return float32(f)
}
