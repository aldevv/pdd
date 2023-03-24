package handlers

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func SendAI(ctx *gin.Context, photoURL string) error {

	// Load the AWS SDK configuration
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == sqs.ServiceID {
			return aws.Endpoint{
				PartitionID: "aws",
				URL:         "http://localhost:9324",
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	// cfg, err := config.LoadDefaultConfig(context.Background(), config.WithEndpointResolver(
	// 	aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
	// 		if service == "sqs" {
	// 			return aws.Endpoint{
	// 				URL: "http://localhost:9324",
	// 			}, nil
	// 		}
	// 		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested: %s/%s", service, region)
	// 	}),
	// ))
	if err != nil {
		return fmt.Errorf("failed to load SDK configuration: %v", err)
	}

	// Create a new SQS service client
	sqsClient := sqs.NewFromConfig(cfg)

	// Set the URL of the SQS queue you want to send the message to
	queueURL := os.Getenv("SQS_QUEUE_URL")

	if queueURL == "" {
		log.Error().Msg("SQS_QUEUE_URL is not defined")
		return errors.New("SQS_QUEUE_URL is not defined")
	}

	// Create the message body by encoding the photo URL as a string
	messageBody := photoURL

	// Send the message to the specified queue
	_, err = sqsClient.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: &messageBody,
	})

	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %v", err)
	}

	fmt.Println("added photoURL to queue successfully")
	return nil
}
