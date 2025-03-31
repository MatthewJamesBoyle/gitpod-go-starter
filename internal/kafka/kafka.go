package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer handles sending messages to Kafka
type Producer struct {
	writer *kafka.Writer
}

// Consumer handles receiving messages from Kafka
type Consumer struct {
	reader *kafka.Reader
}

// Message represents a Kafka message structure
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	Time    time.Time   `json:"time"`
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		BatchTimeout: 10 * time.Millisecond,
	}

	return &Producer{writer: writer}, nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	return p.writer.Close()
}

// SendMessage sends a message to the specified topic
func (p *Producer) SendMessage(ctx context.Context, topic string, msgType string, payload interface{}) error {
	msg := Message{
		Type:    msgType,
		Payload: payload,
		Time:    time.Now(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: data,
		Time:  time.Now(),
	})

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		GroupID:  "go-starter-consumer-group",
	})

	return &Consumer{reader: reader}, nil
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// ConsumeMessages starts consuming messages from Kafka
func (c *Consumer) ConsumeMessages() error {
	ctx := context.Background()
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		var message Message
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Process the message based on its type
		log.Printf("Received message: Type=%s, Time=%v", message.Type, message.Time)

		// Example of handling different message types
		switch message.Type {
		case "todo_created":
			log.Printf("Todo created: %v", message.Payload)
		case "todo_updated":
			log.Printf("Todo updated: %v", message.Payload)
		case "todo_deleted":
			log.Printf("Todo deleted: %v", message.Payload)
		default:
			log.Printf("Unknown message type: %s", message.Type)
		}
	}
}
