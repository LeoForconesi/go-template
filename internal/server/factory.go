package server

import (
	"net/http"

	"github.com/LeonardoForconesi/go-template/internal/server/middleware"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/LeonardoForconesi/go-template/internal/config"
	"github.com/LeonardoForconesi/go-template/pkg/logger"

	kafkagen "github.com/LeonardoForconesi/go-template/internal/adapter/event/kafka"
	rabbitgen "github.com/LeonardoForconesi/go-template/internal/adapter/event/rabbitmq"
	"github.com/LeonardoForconesi/go-template/internal/adapter/repository"

	v1 "github.com/LeonardoForconesi/go-template/internal/adapter/http/v1"
	"github.com/LeonardoForconesi/go-template/internal/service/user"
	useruc "github.com/LeonardoForconesi/go-template/internal/usecase/user"
)

type App struct {
	Log *zap.Logger
	DB  *gorm.DB

	KafkaProducer   *kafkagen.Producer
	RabbitPublisher *rabbitgen.Publisher

	HTTP *http.Server
}

func Build(cfg config.App) (*App, error) {
	log, err := logger.New(cfg.Env)
	if err != nil {
		return nil, err
	}

	db, err := repository.NewPostgres(cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}

	authMW, err := middleware.NewAuthMiddleware(middleware.AuthConfig{
		Domain:   cfg.Auth.Domain,
		Audience: cfg.Auth.Audience,
	})
	if err != nil {
		return nil, err
	}

	// Producers genéricos
	kp, err := kafkagen.NewProducer(log, cfg.Kafka.Brokers)
	if err != nil {
		return nil, err
	}
	rp, err := rabbitgen.NewPublisher(log, cfg.Rabbit.URL)
	if err != nil {
		return nil, err
	}

	// Servicios de dominio
	userEvents := userService.NewUserEventService(log, kp, cfg.Kafka.TopicUsersCreated)
	userNotify := userService.NewUserNotifyService(log, rp, cfg.Rabbit.Queue, cfg.Rabbit.TTLms)

	// Repo
	userRepo := repository.NewUserGormRepository(db)

	// Use cases
	create := &useruc.Creator{Repo: userRepo, Publisher: userEvents}
	get := &useruc.Getter{Repo: userRepo}
	list := &useruc.Lister{Repo: userRepo}
	update := &useruc.Updater{Repo: userRepo}
	del := &useruc.Deleter{Repo: userRepo}
	notify := &useruc.Notifier{Repo: userRepo, Notifier: userNotify}

	handlers := &v1.UserHandlers{
		Create: create,
		Get:    get,
		List:   list,
		Update: update,
		Delete: del,
		Notify: notify,
	}

	router := NewRouter(log, authMW, handlers)

	srv := &http.Server{
		Addr:    cfg.HTTP.Port,
		Handler: router,
	}

	return &App{
		Log:             log,
		DB:              db,
		KafkaProducer:   kp,
		RabbitPublisher: rp,
		HTTP:            srv,
	}, nil
}
