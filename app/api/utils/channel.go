package utils

import (
	"PicSearch/app/db"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func GetEmbeddings(query string) ([]float32, error) {
	var amqpServerURL = os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(amqpServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"generate_clip_encoding", // name
		false,                    // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(query),
		})
	if err != nil {
		return nil, fmt.Errorf("failed to publish a message: %w", err)
	}

	ctx := context.Background()
	sub := db.RedisDB.Subscribe(ctx, "encoding_results")
	defer sub.Close()
	_, err = sub.Receive(ctx)
	if err != nil {
		return nil, err
	}
	msg, err := sub.ReceiveMessage(ctx)

	if err != nil {
		return nil, err
	}
	var date map[string][]float32
	err = json.Unmarshal([]byte(msg.Payload), &date)
	if err != nil {
		fmt.Println("Error deserializing message:", err)
		return nil, err
	}
	return date["embeddings"], nil
}
