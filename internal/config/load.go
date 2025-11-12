package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func Load() (App, error) {
	var cfg App

	v := viper.New()

	// fuentes de configuración
	v.SetConfigName("config")   // busca config.yaml, config.yml, config.env
	v.AddConfigPath(".")        // raíz del proyecto
	v.AddConfigPath("./config") // o dentro de /config
	v.SetConfigType("yaml")     // tipo principal
	v.SetConfigType("env")      // Soporta tambien .env
	v.SetConfigName(".env")
	v.AddConfigPath(".")
	_ = v.MergeInConfig() // mergea si existe

	// Permitir variables de entorno (sobrescriben todo)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // ejemplo: APP.ENV → APP_ENV

	// Defaults
	v.SetDefault("app.env", "dev")
	v.SetDefault("http.port", ":8080")
	v.SetDefault("postgres.dsn", "host=postgres port=5432 user=app password=app dbname=appdb sslmode=disable TimeZone=UTC")
	v.SetDefault("kafka.brokers", []string{"redpanda:9092"})
	v.SetDefault("kafka.topicUsersCreated", "users.created")
	v.SetDefault("rabbit.url", "amqp://guest:guest@rabbitmq:5672/")
	v.SetDefault("rabbit.queue", "user.notify")
	v.SetDefault("rabbit.ttlms", 60000)

	// Intentar leer config.yaml (no falla si no existe)
	if err := v.ReadInConfig(); err == nil {
		log.Printf("config loaded from: %s", v.ConfigFileUsed())
	} else {
		log.Printf("no config.yaml found (using defaults/env): %v", err)
	}

	// Mapear valores al struct principal
	if err := v.Unmarshal(&cfg); err != nil {
		return App{}, fmt.Errorf("config unmarshal: %w", err)
	}

	return cfg, nil
}
