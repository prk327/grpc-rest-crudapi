package config

import "os"

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Schema   string
}

func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnvWithDefault("DB_HOST", "localhost"),
		Port:     getEnvWithDefault("DB_PORT", "5433"),
		User:     getEnvWithDefault("DB_USER", "cognos"),
		Password: getEnvWithDefault("DB_PASSWORD", "admin1234"),
		Name:     getEnvWithDefault("DB_NAME", "vdb"),
		Schema:   getEnvWithDefault("DB_SCHEMA", "omniq"), // Default to "omniq"
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

type ServerConfig struct {
	GRPCPort string
	HTTPPort string
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		GRPCPort: os.Getenv("GRPC_PORT"),
		HTTPPort: os.Getenv("HTTP_PORT"),
	}
}
