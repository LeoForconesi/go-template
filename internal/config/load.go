package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func Load() (App, error) {
	var cfg App

	v := viper.New()
	v.SetConfigType("yaml")

	// Defaults globales (última red de seguridad)
	v.SetDefault("app.env", "local")
	v.SetDefault("http.port", ":8080")
	v.SetDefault("postgres.dsn", "host=localhost port=5432 user=app password=app dbname=appdb sslmode=disable TimeZone=UTC")
	v.SetDefault("kafka.brokers", []string{"localhost:19092"})
	v.SetDefault("kafka.topicUsersCreated", "users.created")
	v.SetDefault("rabbit.url", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("rabbit.queue", "user.notify")
	v.SetDefault("rabbit.ttlms", 60000)

	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		return App{}, fmt.Errorf("load base config.yaml: %w", err)
	}
	log.Printf("base config loaded from: %s", v.ConfigFileUsed())

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	envFile := "config." + env
	v.SetConfigName(envFile)

	if err := v.MergeInConfig(); err != nil {
		log.Printf("no %s.yaml found (using only base + defaults): %v", envFile, err)
	} else {
		log.Printf("env config loaded from: %s.yaml", envFile)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return App{}, fmt.Errorf("config unmarshal: %w", err)
	}

	return cfg, nil
}
