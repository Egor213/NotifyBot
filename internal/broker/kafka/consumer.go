package kafkabroker

import (
	"context"
	"time"

	"github.com/Egor213/notifyBot/internal/broker/worker"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type ConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type KafkaConsumer struct {
	reader  *kafka.Reader
	workers []worker.Worker
}

func NewConsumer(cfg ConsumerConfig) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.Topic,
		GroupID:        cfg.GroupID,
		MinBytes:       1e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		// StartOffset:    kafka.LastOffset,
	})
	return &KafkaConsumer{reader: r}
}

func (c *KafkaConsumer) RegisterWorker(w worker.Worker) {
	c.workers = append(c.workers, w)
}

func (c *KafkaConsumer) Run(ctx context.Context) error {
	log.Debugf("Kafka consumer started for topic '%s'\n", c.reader.Config().Topic)
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Debugf("Kafka consumer stopped.")
				return nil
			}
			log.Debugf("Kafka read error: %v", err)
			continue
		}
		log.Debugf("message: topic=%s partition=%d offset=%d key=%s value=%s\n",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		for _, w := range c.workers {
			go w.ProcessMsg(ctx, m)
		}
	}
}

func (c *KafkaConsumer) Close() error {
	log.Info("Closing Kafka consumer...")
	return c.reader.Close()
}
