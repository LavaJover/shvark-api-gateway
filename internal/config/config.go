package config

import (
	"log"
	"os"
	"github.com/ilyakaznacheev/cleanenv"
)

type HttpAPIConfig struct {
	Env string 	   `yaml:"env"`
	HttpAPIServer  `yaml:"http_server"`
	LogConfig 	   `yaml:"log_config"`
	UserService	   `yaml:"user-service"`
	SSOService	   `yaml:"sso-service"`
	AuthzService   `yaml:"authz-service"`
	OrderService   `yaml:"order-service"`
	BankingService `yaml:"banking-service"`
	WalletService  `yaml:"wallet-service"`
	SwaggerConfig  `yaml:"swagger_config"`
}

type HttpAPIServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type LogConfig struct {
	LogLevel 	string 	`yaml:"log_level"`
	LogFormat 	string 	`yaml:"log_format"`
	LogOutput 	string 	`yaml:"log_output"`
}

type OrderService struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AuthzService struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type UserService struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type SSOService struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type BankingService struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type WalletService struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type SwaggerConfig struct {
	Host 	 string `yaml:"host"`
	Port 	 string `yaml:"port"`
	BasePath string `yaml:"base_path"`
	Schemes  string `yaml:"schemes"`
}

func MustLoad() *HttpAPIConfig {

	// Processing env config variable and file
	configPath := os.Getenv("API_CONFIG_PATH")

	if configPath == ""{
		log.Fatalf("API_CONFIG_PATH was not found\n")
	}

	if _, err := os.Stat(configPath); err != nil{
		log.Fatalf("failed to find config file: %v\n", err)
	}

	// YAML to struct object
	var cfg HttpAPIConfig
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil{
		log.Fatalf("failed to read config file: %v", err)
	}

	return &cfg
}