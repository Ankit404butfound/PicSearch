package utils

import (
	"PicSearch/app/db"
	"PicSearch/app/db/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/streadway/amqp"
)

func TriggerQueue(channelName string, data string) error {

	var amqpServerURL = os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(amqpServerURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		channelName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	return nil
}

func TriggerImageProcessingJob(imageID int) {
	var job models.Job
	job.FileId = imageID
	job.FaceEncodingStatus = "pending"
	job.UniversalEncodingStatus = "pending"

	db.DB.Create(&job)

	TriggerQueue("face_encoder", strconv.Itoa(job.ID))
	TriggerQueue("clip_processor", strconv.Itoa(job.ID))

}

func GetEmbeddings(query string) ([]float32, error) {
	err := TriggerQueue("generate_clip_encoding", query)
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
