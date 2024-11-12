package pubsubservice

import (
	"context"
	"fmt"
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

func ReceiveMessages(subscriptionID string) error {
	if client == nil {
		return fmt.Errorf("Pub/Sub client is not initialized")
	}

	sub := client.Subscription(subscriptionID)

	err := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", string(msg.Data))
		msg.Ack() 
	})
	if err != nil {
		return fmt.Errorf("failed to receive messages: %v", err)
	}

	return nil
}
