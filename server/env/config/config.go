package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type (
	Server struct {
		Mode string `env:"MODE"`
		Port string `env:"PORT"`
	}

	Client struct {
		Url  string `env:"FRONTEND_URL"`
		Port string `env:"CLIENT_PORT"`
	}

	Database struct {
		DBHost     string `env:"DB_HOST"`
		DBPort     string `env:"DB_PORT"`
		DBUser     string `env:"DB_USER"`
		DBPassword string `env:"DB_PASSWORD"`
		DBName     string `env:"DB_NAME"`
	}

	Redis struct {
		RHost string `env:"REDIS_HOST"`
		RPort string `env:"REDIS_PORT"`
	}

	GoogleOAuth struct {
		GOClientID     string `env:"GOOGLE_CLIENT_ID"`
		GOClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	}

	GithubOAuth struct {
		GHClientID     string `env:"GITHUB_CLIENT_ID"`
		GHClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	}

	MicrosoftOAuth struct {
		MSClientID       string `env:"MICROSOFT_CLIENT_ID"`
		MSClientSecret   string `env:"MICROSOFT_CLIENT_SECRET"`
		MSTenantID       string `env:"MICROSOFT_TENANT_ID"`
		MSClientSecretID string `env:"MICROSOFT_CLIENT_SECRET_ID"`
	}

	OAuth struct {
		Google    GoogleOAuth
		Github    GithubOAuth
		Microsoft MicrosoftOAuth
	}

	SMTP struct {
		SHost     string `env:"SMTP_HOST"`
		SPort     string `env:"SMTP_PORT"`
		SUser     string `env:"SMTP_USER"`
		SPassword string `env:"SMTP_PASSWORD"`
	}

	RabbitMQ struct {
		RMQHost        string `env:"RABBITMQ_HOST"`
		RMQPort        string `env:"RABBITMQ_PORT"`
		RMQUser        string `env:"RABBITMQ_USER"`
		RMQPassword    string `env:"RABBITMQ_PASSWORD"`
		RMQVirtualHost string `env:"RABBITMQ_VIRTUAL_HOST"`
	}

	Config struct {
		Server   Server
		Client   Client
		Database Database
		Redis    Redis
		OAuth    OAuth
		SMTP     SMTP
		RabbitMQ RabbitMQ
	}
)

var Cfg Config

func Load() {
	var ok bool

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Error loading .env file")
		}
	}

	// ! Load Server configuration ____________________________
	if Cfg.Server.Mode, ok = os.LookupEnv("MODE"); !ok {
		log.Println("MODE env is not set")
	}
	if Cfg.Server.Port, ok = os.LookupEnv("PORT"); !ok {
		log.Println("PORT env is not set")
	}
	// ! ______________________________________________________

	// ! Load Client configuration ____________________________
	if Cfg.Client.Url, ok = os.LookupEnv("FRONTEND_URL"); !ok {
		log.Println("FRONTEND_URL env is not set")
	}
	if Cfg.Client.Port, ok = os.LookupEnv("CLIENT_PORT"); !ok {
		log.Println("CLIENT_PORT env is not set")
	}
	// ! ______________________________________________________

	// ! Load Database configuration __________________________
	if Cfg.Database.DBUser, ok = os.LookupEnv("DB_USER"); !ok {
		log.Println("DB_USER env is not set")
	}
	if Cfg.Database.DBHost, ok = os.LookupEnv("DB_HOST"); !ok {
		log.Println("DB_HOST env is not set")
	}
	if Cfg.Database.DBPort, ok = os.LookupEnv("DB_PORT"); !ok {
		log.Println("DB_PORT env is not set")
	}
	if Cfg.Database.DBName, ok = os.LookupEnv("DB_NAME"); !ok {
		log.Println("DB_NAME env is not set")
	}
	if Cfg.Database.DBPassword, ok = os.LookupEnv("DB_PASSWORD"); !ok {
		log.Println("DB_PASSWORD env is not set")
	}
	// ! ______________________________________________________

	// ! Load Redis configuration _____________________________
	if Cfg.Redis.RHost, ok = os.LookupEnv("REDIS_HOST"); !ok {
		log.Println("REDIS_HOST env is not set")
	}
	if Cfg.Redis.RPort, ok = os.LookupEnv("REDIS_PORT"); !ok {
		log.Println("REDIS_PORT env is not set")
	}
	// ! ______________________________________________________

	// ! Load Google OAuth configuration ______________________
	if Cfg.OAuth.Google.GOClientID, ok = os.LookupEnv("GOOGLE_CLIENT_ID"); !ok {
		log.Println("GOOGLE_CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Google.GOClientSecret, ok = os.LookupEnv("GOOGLE_CLIENT_SECRET"); !ok {
		log.Println("GOOGLE_CLIENT_SECRET env is not set")
	}
	// ! ______________________________________________________

	// ! Load Github OAuth configuration ______________________
	if Cfg.OAuth.Github.GHClientID, ok = os.LookupEnv("GITHUB_CLIENT_ID"); !ok {
		log.Println("GITHUB_CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Github.GHClientSecret, ok = os.LookupEnv("GITHUB_CLIENT_SECRET"); !ok {
		log.Println("GITHUB_CLIENT_SECRET env is not set")
	}
	// ! ______________________________________________________

	// ! Load Microsoft OAuth configuration ___________________
	if Cfg.OAuth.Microsoft.MSClientID, ok = os.LookupEnv("MICROSOFT_CLIENT_ID"); !ok {
		log.Println("MICROSOFT_CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Microsoft.MSClientSecret, ok = os.LookupEnv("MICROSOFT_CLIENT_SECRET"); !ok {
		log.Println("MICROSOFT_CLIENT_SECRET env is not set")
	}
	if Cfg.OAuth.Microsoft.MSTenantID, ok = os.LookupEnv("MICROSOFT_TENANT_ID"); !ok {
		log.Println("MICROSOFT_TENANT_ID env is not set")
	}
	if Cfg.OAuth.Microsoft.MSClientSecretID, ok = os.LookupEnv("MICROSOFT_CLIENT_SECRET_ID"); !ok {
		log.Println("MICROSOFT_CLIENT_SECRET_ID env is not set")
	}
	// ! ______________________________________________________

	// ! Load SMTP configuration ______________________________
	if Cfg.SMTP.SHost, ok = os.LookupEnv("SMTP_HOST"); !ok {
		log.Println("SMTP_HOST env is not set")
	}
	if Cfg.SMTP.SPort, ok = os.LookupEnv("SMTP_PORT"); !ok {
		log.Println("SMTP_PORT env is not set")
	}
	if Cfg.SMTP.SUser, ok = os.LookupEnv("SMTP_USER"); !ok {
		log.Println("SMTP_USER env is not set")
	}
	if Cfg.SMTP.SPassword, ok = os.LookupEnv("SMTP_PASSWORD"); !ok {
		log.Println("SMTP_PASSWORD env is not set")
	}
	// ! ______________________________________________________

	// ! Load RabbitMQ configuration __________________________
	if Cfg.RabbitMQ.RMQUser, ok = os.LookupEnv("RABBITMQ_USER"); !ok {
		log.Println("RABBITMQ_USER env is not set")
	}
	if Cfg.RabbitMQ.RMQPassword, ok = os.LookupEnv("RABBITMQ_PASSWORD"); !ok {
		log.Println("RABBITMQ_PASSWORD env is not set")
	}
	if Cfg.RabbitMQ.RMQHost, ok = os.LookupEnv("RABBITMQ_HOST"); !ok {
		log.Println("RABBITMQ_HOST env is not set")
	}
	if Cfg.RabbitMQ.RMQPort, ok = os.LookupEnv("RABBITMQ_PORT"); !ok {
		log.Println("RABBITMQ_PORT env is not set")
	}
	if Cfg.RabbitMQ.RMQVirtualHost, ok = os.LookupEnv("RABBITMQ_VIRTUAL_HOST"); !ok {
		log.Println("RABBITMQ_VIRTUAL_HOST env is not set")
	}
	// ! ______________________________________________________
}
