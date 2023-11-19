package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"PORT"`
	DBUrl        string `mapstructure:"DB_URL"`
	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
}

func SetDefaultConfig() {
	viper.SetDefault("PORT", ":50051")
	viper.SetDefault("DB_URL", "postgres://postgres:postgres@127.0.0.1:5432/auth_svc")
	viper.SetDefault("JWT_SECRET_KEY", "supersecret45632*/*-")
}

func ReadSystemEnv() {

	env_port, ok := os.LookupEnv("PORT")
	if !ok || env_port == "" {
		fmt.Println("System environment variables 'PORT' not set, Working with default values instead...")
		viper.SetDefault("PORT", ":50051")
	}
	viper.BindEnv("PORT")

	env_db_url, ok := os.LookupEnv("DB_URL")
	if !ok || env_db_url == "" {
		fmt.Println("System environment variables 'DB_URL' not set, Working with default values instead...")
		viper.SetDefault("DB_URL", "postgres://postgres:postgres@127.0.0.1:5432/auth_svc")
	}
	viper.BindEnv("DB_URL")

	env_jwt, ok := os.LookupEnv("JWT_SECRET_KEY")
	if !ok || env_jwt == "" {
		fmt.Println("System environment variables 'JWT_SECRET_KEY' not set, Working with default values instead...")
		viper.SetDefault("JWT_SECRET_KEY", ":supersecret45632*/*-")
	}
	viper.BindEnv("JWT_SECRET_KEY")

}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No configuration file found, checking system environment variables instead ...")
			ReadSystemEnv()
		} else {
			fmt.Println("Error proceeded while trying to set environment variables, check your configuration")
			fmt.Println("Working with default values instead...")
			SetDefaultConfig()
		}

	}

	err = viper.Unmarshal(&config)
	fmt.Println(os.LookupEnv("DB_URL"))
	fmt.Println(config)

	return
}
