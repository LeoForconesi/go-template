package config

type Postgres struct {
	DSN string `env:"POSTGRES_DSN,required"`
}

type Kafka struct {
	Brokers           []string `env:"KAFKA_BROKERS,required" envSeparator:","`
	TopicUsersCreated string   `env:"KAFKA_TOPIC_USERS_CREATED,required" default:"users.created"`
}

type Rabbit struct {
	URL   string `env:"RABBIT_URL,required"`   // "amqp://guest:guest@localhost:5672/"
	Queue string `env:"RABBIT_QUEUE,required"` // "user.notify"
	TTLms int    `env:"RABBIT_TTL_MS" envDefault:"60000"`
}

type HTTP struct {
	Port string `env:"HTTP_PORT" envDefault:":8080"`
}

type Auth struct {
	Domain   string `env:"AUTH_DOMAIN" envDefault:"https://example.com"`
	Audience string `env:"AUTH_AUDIENCE" envDefault:"https://example.com/api"`
}

type App struct {
	Env      string `env:"APP_ENV" envDefault:"dev"` // dev|prod|test
	Postgres Postgres
	Kafka    Kafka
	Rabbit   Rabbit
	HTTP     HTTP
	Auth     Auth
}
