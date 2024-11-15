package pubsubservice

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

var (
	client *pubsub.Client
)

func InitializePubSubClient(projectID, credentialsPath string) error {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentialsPath)

	ctx := context.Background()
	var err error
	client, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to create Pub/Sub client: %v", err)
	}

	return nil
}

func PublishMessage(topicID, msg string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("Pub/Sub client is not initialized")
	}

	topic := client.Topic(topicID)
	result := topic.Publish(context.Background(), &pubsub.Message{
		Data: []byte(msg),
	})

	id, err := result.Get(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to publish message: %v", err)
	}

	return id, nil
}

func ReceiveMessages(subscriptionID string, updateChan chan<- string) {
	if client == nil {
		log.Println("Pub/Sub client is not initialized")
		close(updateChan)
		return
	}

	sub := client.Subscription(subscriptionID)

	go func() {
		err := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			updateChan <- string(msg.Data)
			msg.Ack()
		})

		if err != nil {
			log.Printf("Error receiving messages: %v", err)
		}
	}()
}
