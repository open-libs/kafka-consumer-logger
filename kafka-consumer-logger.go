package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var logger *zap.Logger

var latest_message kafka.Message

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	return cfg.Build()
}

func main_handler(w http.ResponseWriter, r *http.Request) {
	result, _ := json.Marshal(latest_message)
	w.Write(result)
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6, // 10MB
		MaxWait:  time.Second,
	})
}

func run_consumer() {
	// get kafka reader using environment variables.
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		logger.Error("environment variable not found: KAFKA_URL")
		os.Exit(1)
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		logger.Error("environment variable not found: KAFKA_TOPIC")
		os.Exit(1)
	}
	groupID := os.Getenv("KAFKA_GROUP_ID")

	if groupID == "" {
		logger.Error("environment variable not found: KAFKA_GROUP_ID")
		os.Exit(1)
	}

	reader := getKafkaReader(kafkaURL, topic, groupID)

	defer reader.Close()

	logger.Info("start consuming",
		zap.String("KAFKA_URL", kafkaURL),
		zap.String("KAFKA_TOPIC", topic),
		zap.String("KAFKA_GROUP_ID", groupID),
	)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		latest_message = m
		logger.Info("message received",
			zap.String("topic", m.Topic),
			zap.Int("partition", m.Partition),
			zap.Int64("offset", m.Offset),
			zap.String("key", string(m.Key)),
			zap.String("value", string(m.Value)),
		)
	}
}

func main() {
	logger, _ = NewLogger()
	go run_consumer()
	http.HandleFunc("/", main_handler)
	listen := os.Getenv("HTTP_PORT")
	if listen == "" {
		listen = ":80"
	}
	logger.Info("Listening on " + listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}
