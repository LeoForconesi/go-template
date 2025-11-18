package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

/*
Sarama
Qué es: una librería Go para hablar con Apache Kafka (producir y consumir mensajes).
Qué hace: te da productores/consumidores, control de acks, compresión, retries, particionado, idempotencia, etc.
Alternativas: segmentio/kafka-go (API simple) o drivers nativos de Confluent (más “enterprise”).
*/

type Producer struct {
	log *zap.Logger
	sp  sarama.SyncProducer
}

func NewProducer(log *zap.Logger, brokers []string) (*Producer, error) {
	cfg := sarama.NewConfig()

	// Requisitos para idempotencia
	cfg.Version = sarama.V3_5_0_0
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1
	cfg.Producer.Retry.Max = 5

	sp, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &Producer{log: log, sp: sp}, nil
}

// Produce envía un mensaje genérico a un tópico Kafka.
func (p *Producer) Produce(_ context.Context, topic string, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.ByteEncoder(key),
		Value:     sarama.ByteEncoder(value),
		Timestamp: time.Now(),
	}
	_, _, err := p.sp.SendMessage(msg)
	if err != nil {
		p.log.Warn("kafka publish failed", zap.Error(err), zap.String("topic", topic))
		return err
	}
	p.log.Debug("kafka message produced", zap.String("topic", topic))
	return nil
}

func (p *Producer) Close() error {
	return p.sp.Close()
}
