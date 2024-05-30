package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AwsConfig struct {
	Region            string `mapstructure:"region"`
	CognitoClientId   string `mapstructure:"cognito_client_id"`
	CognitoUserPoolID string `mapstructure:"cognito_user_pool_id"`
}

type ApiConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type SQLDatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type Config struct {
	Aws AwsConfig         `mapstructure:"aws"`
	Api ApiConfig         `mapstructure:"api"`
	Sql SQLDatabaseConfig `mapstructure:"sql"`
	Env string            `mapstructure:"env"`
}

func setDefaults() {
	viper.SetDefault("env", "development")

	viper.SetDefault("aws.region", "us-east-1")
	viper.SetDefault("aws.cognito_client_id", "SET_ME")
	viper.SetDefault("aws.cognito_user_pool_id", "SET_ME")

	viper.SetDefault("sql.host", "localhost")
	viper.SetDefault("sql.port", 5432)
	viper.SetDefault("sql.user", "SET_ME")
	viper.SetDefault("sql.password", "SET_ME")
	viper.SetDefault("sql.database", "SET_ME")

	viper.SetDefault("api.host", "0.0.0.0")
	viper.SetDefault("api.port", 4000)
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfigAs(configPath + "/config.yaml"); err != nil {
				return nil, fmt.Errorf("error writing default config file: %v", err)
			}
		} else {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	return &config, nil
}
