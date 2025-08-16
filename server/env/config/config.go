package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Server struct {
		Mode         string `env:"MODE"`
		Port         string `env:"PORT"`
		JWTSecretKey string `env:"JWT_SECRET_KEY"`
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

	GSMTP struct {
		GSHost     string `env:"GOOGLE_SMTP_HOST"`
		GSPort     string `env:"GOOGLE_SMTP_PORT"`
		GSUser     string `env:"GOOGLE_SMTP_USER"`
		GSPassword string `env:"GOOGLE_SMTP_PASSWORD"`
	}

	RabbitMQ struct {
		RMQHost        string `env:"RABBITMQ_HOST"`
		RMQPort        string `env:"RABBITMQ_PORT"`
		RMQUser        string `env:"RABBITMQ_USER"`
		RMQPassword    string `env:"RABBITMQ_PASSWORD"`
		RMQVirtualHost string `env:"RABBITMQ_VIRTUAL_HOST"`
	}

	ZSMTP struct {
		ZSHost     string `env:"ZOHO_SMTP_HOST"`
		ZSPort     string `env:"ZOHO_SMTP_PORT"`
		ZSUser     string `env:"ZOHO_SMTP_USER"`
		ZSPassword string `env:"ZOHO_SMTP_PASSWORD"`
		ZSSecure   string `env:"ZOHO_SMTP_SECURE"`
		ZSAuth     bool   `env:"ZOHO_SMTP_AUTH"`
	}

	Config struct {
		Server   Server
		Client   Client
		Database Database
		Redis    Redis
		OAuth    OAuth
		GSMTP    GSMTP
		ZSMTP    ZSMTP
		RabbitMQ RabbitMQ
	}
)

var Cfg Config

func LoadNative() {
	var ok bool

	envFile := os.Getenv("APP_ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("[ERROR] No %s file found\n", envFile)
	}

	// ! Load Server configuration ____________________________
	if Cfg.Server.Mode, ok = os.LookupEnv("MODE"); !ok {
		log.Println("MODE env is not set")
	}
	if Cfg.Server.Port, ok = os.LookupEnv("PORT"); !ok {
		log.Println("PORT env is not set")
	}
	if Cfg.Server.JWTSecretKey, ok = os.LookupEnv("JWT_SECRET_KEY"); !ok {
		log.Println("JWT_SECRET_KEY env is not set")
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

	// ! Load Gmail SMTP configuration ______________________________
	if Cfg.GSMTP.GSHost, ok = os.LookupEnv("GOOGLE_SMTP_HOST"); !ok {
		log.Println("GOOGLE_SMTP_HOST env is not set")
	}
	if Cfg.GSMTP.GSPort, ok = os.LookupEnv("GOOGLE_SMTP_PORT"); !ok {
		log.Println("GOOGLE_SMTP_PORT env is not set")
	}
	if Cfg.GSMTP.GSUser, ok = os.LookupEnv("GOOGLE_SMTP_USER"); !ok {
		log.Println("GOOGLE_SMTP_USER env is not set")
	}
	if Cfg.GSMTP.GSPassword, ok = os.LookupEnv("GOOGLE_SMTP_PASSWORD"); !ok {
		log.Println("GOOGLE_SMTP_PASSWORD env is not set")
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

	// ! Load Zoho SMTP configuration __________________________
	if Cfg.ZSMTP.ZSHost, ok = os.LookupEnv("ZOHO_SMTP_HOST"); !ok {
		log.Println("ZOHO_SMTP_HOST env is not set")
	}
	if Cfg.ZSMTP.ZSPort, ok = os.LookupEnv("ZOHO_SMTP_PORT"); !ok {
		log.Println("ZOHO_SMTP_PORT env is not set")
	}
	if Cfg.ZSMTP.ZSUser, ok = os.LookupEnv("ZOHO_SMTP_USER"); !ok {
		log.Println("ZOHO_SMTP_USER env is not set")
	}
	if Cfg.ZSMTP.ZSPassword, ok = os.LookupEnv("ZOHO_SMTP_PASSWORD"); !ok {
		log.Println("ZOHO_SMTP_PASSWORD env is not set")
	}
	if Cfg.ZSMTP.ZSSecure, ok = os.LookupEnv("ZOHO_SMTP_SECURE"); !ok {
		log.Println("ZOHO_SMTP_SECURE env is not set")
	}
	if zohoAuth, ok := os.LookupEnv("ZOHO_SMTP_AUTH"); !ok {
		log.Println("ZOHO_SMTP_AUTH env is not set")
	} else {
		Cfg.ZSMTP.ZSAuth = zohoAuth == "true"
	}
	// ! ______________________________________________________
}

func LoadByViper() error {
	config := viper.New()
	config.SetConfigFile("app/config.json")

	if err := config.ReadInConfig(); err != nil {
		return err
	}

	// ! Load Server configuration ____________________________
	if Cfg.Server.Mode = config.GetString("MODE"); Cfg.Server.Mode == "" {
		log.Println("MODE env is not set")
	}
	if Cfg.Server.Port = config.GetString("PORT"); Cfg.Server.Port == "" {
		log.Println("PORT env is not set")
	}
	if Cfg.Server.JWTSecretKey = config.GetString("JWT_SECRET_KEY"); Cfg.Server.JWTSecretKey == "" {
		log.Println("JWT_SECRET_KEY env is not set")
	}
	// ! ______________________________________________________

	// ! Load Client configuration ____________________________
	if Cfg.Client.Url = config.GetString("CLIENT.URL"); Cfg.Client.Url == "" {
		log.Println("CLIENT.URL env is not set")
	}
	if Cfg.Client.Port = config.GetString("CLIENT.PORT"); Cfg.Client.Port == "" {
		log.Println("CLIENT.PORT env is not set")
	}
	// ! ______________________________________________________

	// ! Load Database configuration __________________________
	if Cfg.Database.DBUser = config.GetString("DATABASE.POSTGRESQL.USER"); Cfg.Database.DBUser == "" {
		log.Println("DATABASE.POSTGRESQL.USER env is not set")
	}
	if Cfg.Database.DBHost = config.GetString("DATABASE.POSTGRESQL.HOST"); Cfg.Database.DBHost == "" {
		log.Println("DATABASE.POSTGRESQL.HOST env is not set")
	}
	if Cfg.Database.DBPort = config.GetString("DATABASE.POSTGRESQL.PORT"); Cfg.Database.DBPort == "" {
		log.Println("DATABASE.POSTGRESQL.PORT env is not set")
	}
	if Cfg.Database.DBName = config.GetString("DATABASE.POSTGRESQL.NAME"); Cfg.Database.DBName == "" {
		log.Println("DATABASE.POSTGRESQL.NAME env is not set")
	}
	if Cfg.Database.DBPassword = config.GetString("DATABASE.POSTGRESQL.PASSWORD"); Cfg.Database.DBPassword == "" {
		log.Println("DATABASE.POSTGRESQL.PASSWORD env is not set")
	}
	// ! ______________________________________________________

	// ! Load Redis configuration _____________________________
	if Cfg.Redis.RHost = config.GetString("REDIS.HOST"); Cfg.Redis.RHost == "" {
		log.Println("REDIS.HOST env is not set")
	}
	if Cfg.Redis.RPort = config.GetString("REDIS.PORT"); Cfg.Redis.RPort == "" {
		log.Println("REDIS.PORT env is not set")
	}
	// ! ______________________________________________________

	// ! Load RabbitMQ configuration __________________________
	if Cfg.RabbitMQ.RMQUser = config.GetString("MESSAGE-BROKER.RABBITMQ.USER"); Cfg.RabbitMQ.RMQUser == "" {
		log.Println("MESSAGE-BROKER.RABBITMQ.USER env is not set")
	}
	if Cfg.RabbitMQ.RMQPassword = config.GetString("MESSAGE-BROKER.RABBITMQ.PASSWORD"); Cfg.RabbitMQ.RMQPassword == "" {
		log.Println("MESSAGE-BROKER.RABBITMQ.PASSWORD env is not set")
	}
	if Cfg.RabbitMQ.RMQHost = config.GetString("MESSAGE-BROKER.RABBITMQ.HOST"); Cfg.RabbitMQ.RMQHost == "" {
		log.Println("MESSAGE-BROKER.RABBITMQ.HOST env is not set")
	}
	if Cfg.RabbitMQ.RMQPort = config.GetString("MESSAGE-BROKER.RABBITMQ.PORT"); Cfg.RabbitMQ.RMQPort == "" {
		log.Println("MESSAGE-BROKER.RABBITMQ.PORT env is not set")
	}
	if Cfg.RabbitMQ.RMQVirtualHost = config.GetString("MESSAGE-BROKER.RABBITMQ.VIRTUAL_HOST"); Cfg.RabbitMQ.RMQVirtualHost == "" {
		log.Println("MESSAGE-BROKER.RABBITMQ.VIRTUAL_HOST env is not set")
	}
	// ! ______________________________________________________

	// ! Load Google OAuth configuration ______________________
	if Cfg.OAuth.Google.GOClientID = config.GetString("OAUTH.GOOGLE.CLIENT_ID"); Cfg.OAuth.Google.GOClientID == "" {
		log.Println("OAUTH.GOOGLE.CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Google.GOClientSecret = config.GetString("OAUTH.GOOGLE.CLIENT_SECRET"); Cfg.OAuth.Google.GOClientSecret == "" {
		log.Println("OAUTH.GOOGLE.CLIENT_SECRET env is not set")
	}
	// ! ______________________________________________________

	// ! Load Github OAuth configuration ______________________
	if Cfg.OAuth.Github.GHClientID = config.GetString("OAUTH.GITHUB.CLIENT_ID"); Cfg.OAuth.Github.GHClientID == "" {
		log.Println("OAUTH.GITHUB.CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Github.GHClientSecret = config.GetString("OAUTH.GITHUB.CLIENT_SECRET"); Cfg.OAuth.Github.GHClientSecret == "" {
		log.Println("OAUTH.GITHUB.CLIENT_SECRET env is not set")
	}
	// ! ______________________________________________________

	// ! Load Microsoft OAuth configuration ___________________
	if Cfg.OAuth.Microsoft.MSClientID = config.GetString("OAUTH.MICROSOFT.CLIENT_ID"); Cfg.OAuth.Microsoft.MSClientID == "" {
		log.Println("OAUTH.MICROSOFT.CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Microsoft.MSClientSecret = config.GetString("OAUTH.MICROSOFT.CLIENT_SECRET"); Cfg.OAuth.Microsoft.MSClientSecret == "" {
		log.Println("OAUTH.MICROSOFT.CLIENT_SECRET env is not set")
	}
	if Cfg.OAuth.Microsoft.MSTenantID = config.GetString("OAUTH.MICROSOFT.TENANT_ID"); Cfg.OAuth.Microsoft.MSTenantID == "" {
		log.Println("OAUTH.MICROSOFT.TENANT_ID env is not set")
	}
	if Cfg.OAuth.Microsoft.MSClientSecretID = config.GetString("OAUTH.MICROSOFT.CLIENT_SECRET_ID"); Cfg.OAuth.Microsoft.MSClientSecretID == "" {
		log.Println("OAUTH.MICROSOFT.CLIENT_SECRET_ID env is not set")
	}
	// ! ______________________________________________________

	// ! Load Gmail SMTP configuration ______________________________
	if Cfg.GSMTP.GSHost = config.GetString("SMTP.GOOGLE.HOST"); Cfg.GSMTP.GSHost == "" {
		log.Println("SMTP.GOOGLE.HOST env is not set")
	}
	if Cfg.GSMTP.GSPort = config.GetString("SMTP.GOOGLE.PORT"); Cfg.GSMTP.GSPort == "" {
		log.Println("SMTP.GOOGLE.PORT env is not set")
	}
	if Cfg.GSMTP.GSUser = config.GetString("SMTP.GOOGLE.USER"); Cfg.GSMTP.GSUser == "" {
		log.Println("SMTP.GOOGLE.USER env is not set")
	}
	if Cfg.GSMTP.GSPassword = config.GetString("SMTP.GOOGLE.PASSWORD"); Cfg.GSMTP.GSPassword == "" {
		log.Println("SMTP.GOOGLE.PASSWORD env is not set")
	}
	// ! ______________________________________________________

	// ! Load Zoho SMTP configuration __________________________
	if Cfg.ZSMTP.ZSHost = config.GetString("SMTP.ZOHO.HOST"); Cfg.ZSMTP.ZSHost == "" {
		log.Println("SMTP.ZOHO.HOST env is not set")
	}
	if Cfg.ZSMTP.ZSPort = config.GetString("SMTP.ZOHO.PORT"); Cfg.ZSMTP.ZSPort == "" {
		log.Println("SMTP.ZOHO.PORT env is not set")
	}
	if Cfg.ZSMTP.ZSUser = config.GetString("SMTP.ZOHO.USER"); Cfg.ZSMTP.ZSUser == "" {
		log.Println("SMTP.ZOHO.USER env is not set")
	}
	if Cfg.ZSMTP.ZSPassword = config.GetString("SMTP.ZOHO.PASSWORD"); Cfg.ZSMTP.ZSPassword == "" {
		log.Println("SMTP.ZOHO.PASSWORD env is not set")
	}
	if Cfg.ZSMTP.ZSSecure = config.GetString("SMTP.ZOHO.SECURE"); Cfg.ZSMTP.ZSSecure == "" {
		log.Println("SMTP.ZOHO.SECURE env is not set")
	}
	if Cfg.ZSMTP.ZSAuth = config.GetBool("SMTP.ZOHO.AUTH"); !Cfg.ZSMTP.ZSAuth {
		log.Println("SMTP.ZOHO.AUTH env is not set")
	}
	// ! ______________________________________________________

	return nil
}
