package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Port     string `env:"SERVER_PORT" env-default:"7777" env-description:"Порт сервера"`
		Host     string `env:"SERVER_HOST" env-default:"localhost" env-description:"Хост сервера"`
		LogLevel string `env:"SERVER_LEVEL_LOG" env-default:"info" env-description:"Уровень логирования"`
	}

	Client struct {
		Port string `env:"CLIENT_PORT" env-default:"7777" env-description:"Порт стороннего хранилища"`
		Host string `env:"CLIENT_HOST" env-default:"localhost" env-description:"Хост стороннего хранилища"`
	}

	Database struct {
		Host                   string `env:"DB_HOST" env-default:"localhost" env-description:"Хост базы данных"`
		Port                   string `env:"DB_PORT" env-default:"7777" env-description:"Порт базы данных"`
		Type                   string `env:"DB_TYPE" env-default:"postgres" env-description:"Тип БД"`
		User                   string `env:"DB_USER" env-default:"postgres" env-description:"Пользователь БД"`
		Password               string `env:"DB_PASSWORD" env-default:"postgres" env-description:"Пароль к БД"`
		Name                   string `env:"DB_NAME" env-default:"postgres" env-description:"Наименование БД"`
		SSLMode                string `env:"DB_SSLMODE" env-default:"disable" env-description:"Используем SSL?"`
		MaxOpenConns           int    `env:"DB_MAX_OPEN_CONNS" env-default:"25" env-description:"Макс. кол-во открытых соединений к БД"`
		MaxIdleConns           int    `env:"DB_MAX_IDLE_CONNS" env-default:"25" env-description:"Макс. количество неактивных соединений к БД"`
		ConnMaxIdleTime        string `env:"DB_CONN_MAX_IDLE_TIME" env-default:"5m" env-description:"Макс. время ожидания для неактивных соединений в пуле БД"`
		ConnMaxLifetime        string `env:"DB_CONN_MAX_LIFETIME" env-default:"1h" env-description:"Макс. время жизни пула соединений с БД"`
		ConnectTimeout         int    `env:"DB_CONNECT_TIMEOUT" env-default:"10" env-description:"Время ожидания подключения к БД, [с]"`
		StatementTimeout       int    `env:"DB_STATEMENT_TIMEOUT" env-default:"5000" env-description:"Макс время выполнения SQL-запроса, [мс]"`                                             //Резерв
		IdleInTxSessionTimeout int    `env:"DB_IDLE_IN_TRANSACTION_SESSION_TIMEOUT" env-default:"30000" env-description:"Таймаут для завершения транзакции, если сессия простаивает, [мс]"` //Резерв
	}
}

const fNameConfig = ".env"

const (
	isntExistFileCfg = "файл конфигурации не существует: %s"
	errLoadConfig    = "не смогли прочитат конфигурацию: %s"
)

func MustLoad() *Config {
	workingDirApp, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fConfigPath := filepath.Join(workingDirApp, fNameConfig)

	var cfg Config
	err = cleanenv.ReadConfig(fConfigPath, &cfg)
	if err != nil {
		log.Fatalf(errLoadConfig, err)
	}
	log.Println("Конфигурация загружена")
	return &cfg
}
