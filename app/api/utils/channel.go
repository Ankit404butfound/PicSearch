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
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(query),
	}

	if err := channelRabbitMQ.Publish(
		"",                       // exchange
		"generate_clip_encoding", // queue name
		false,                    // mandatory
		false,                    // immediate
		message,                  // message to publish
	); err != nil {
		return nil, err
	}

	fmt.Print(query)
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

	var floats []float32
	err = json.Unmarshal([]byte(msg.Payload), &floats)
	if err != nil {
		fmt.Println("Error deserializing message:", err)
		return nil, err
	}
	return floats, nil
}
