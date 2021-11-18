package services

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

func PublishTopicMessages(ctx context.Context, projectId, topicId string, messages [][]byte) error {
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	t := client.Topic(topicId)

	var publishedResults []*pubsub.PublishResult

	for _, msg := range messages {
		result := t.Publish(ctx, &pubsub.Message{Data: msg})
		publishedResults = append(publishedResults, result)
	}

	var pubsubErrs []error

	for _, result := range publishedResults {
		id, err := result.Get(ctx)
		if err != nil {
			pubsubErrs = append(pubsubErrs, err)
		}
		log.Printf("Published message; msg ID: %s", id)
	}

	if len(pubsubErrs) != 0 {
		return fmt.Errorf("%d out of %d failed while publishing messages to PubSub\n%v", len(pubsubErrs), len(publishedResults), pubsubErrs)
	}
	return nil
}
