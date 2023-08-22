package config

import (
	"flag"
	"net"
	"os"
)

type Config struct {
	StoragePath string
	ServerAddr  string
	BaseURL     string
	DataBaseDSN string
	TypeStorage string
}

type ConfigBuilder struct {
	config *Config
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{},
	}
}

func (b *ConfigBuilder) Storage(storagePath string) *ConfigBuilder {
	b.config.StoragePath = storagePath
	return b
}

func (b *ConfigBuilder) Address(serverAddr string) *ConfigBuilder {
	b.config.ServerAddr = serverAddr
	return b
}

func (b *ConfigBuilder) BaseURL(baseURL string) *ConfigBuilder {
	b.config.BaseURL = baseURL
	return b
}

func (b *ConfigBuilder) DataBase(dataBase string) *ConfigBuilder {
	b.config.DataBaseDSN = dataBase
	return b
}

func (b *ConfigBuilder) TypeStorage(TypeStorage string) *ConfigBuilder {
	b.config.TypeStorage = TypeStorage
	return b
}

func (b *ConfigBuilder) Build() *Config {
	return b.config
}

func getEnvOrFlag(envKey string, flagValue string, defaultValue string) string {
	envVal := os.Getenv(envKey)
	if envVal == "" && flagValue != "" {
		os.Setenv(envKey, flagValue)
		return flagValue
	} else if envVal == "" {
		os.Setenv(envKey, defaultValue)
		return defaultValue
	}
	return envVal
}

func InitConfig() (*Config, error) {
	var (
		addrFlag     string
		baseURLFlag  string
		fileFlag     string
		dataBaseFlag string
		typeStor     string
	)

	flag.StringVar(&addrFlag, "a", "", "HTTP-сервера")
	flag.StringVar(&baseURLFlag, "b", "", "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&fileFlag, "f", "", "Путь до файла с сокращёнными URL")
	flag.StringVar(&dataBaseFlag, "d", "", "Подключение к базе данных")
	flag.Parse()

	serverAddress := getEnvOrFlag("SERVER_ADDRESS", addrFlag, "127.0.0.1:8080")
	baseURL := getEnvOrFlag("BASE_URL", baseURLFlag, "http://127.0.0.1:8080")
	fileStorage := getEnvOrFlag("FILE_STORAGE_PATH", fileFlag, "./")
	dataBaseDsn := getEnvOrFlag("DATABASE_DSN", dataBaseFlag, "")

	_, err := net.ResolveTCPAddr("tcp", serverAddress)
	if err != nil {
		serverAddress = "127.0.0.1:8080"
		baseURL = "http://127.0.0.1:8080"
	}

	if dataBaseDsn == "" {
		if fileStorage == "./"{
			typeStor = "in-memory"
		} else {
			typeStor = "file"
		}
	} else {
		typeStor = "file"
	}

	os.Setenv("SERVER_ADDRESS", serverAddress)
	os.Setenv("BASE_URL", baseURL)
	os.Setenv("FILE_STORAGE_PATH", fileStorage)
	os.Setenv("DATABASE_DSN", dataBaseDsn)

	builder := NewConfigBuilder().
		Address(serverAddress).
		BaseURL(baseURL).
		Storage(fileStorage).
		DataBase(dataBaseDsn).
		TypeStorage(typeStor)

	return builder.Build(), nil
}
