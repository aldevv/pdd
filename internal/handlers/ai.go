package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/db"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
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
	authToken, exists := ctx.Get("authtoken")
	if !exists {
		log.Error().Msg("not authenticated")
		return errors.New("not authenticated")
	}

	var m = struct {
		PhotoURL  string `json:"photoURL"`
		AuthToken string `json:"authToken"`
	}{PhotoURL: photoURL, AuthToken: authToken.(string)}
	messageBody, err := json.Marshal(m)
	if err != nil {
		log.Error().Msg("message format incorrect")
		return err
	}

	// Send the message to the specified queue
	_, err = sqsClient.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(messageBody)),
	})

	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %v", err)
	}

	fmt.Println("added photoURL to queue successfully")
	return nil
}

type Result struct {
	PhotoURL string `json:"photo_url"`
	Sickness string `json:"sickness"`
	Accuracy string `json:"accuracy"`
}

type Query struct {
	PhotoURL string `json:"photo_url"`
}

func GetResult(ctx *gin.Context) {
	var res Query
	if err := ctx.BindJSON(&res); err != nil {
		fmt.Println("error decoding body: " + err.Error())
		ctx.Status(400)
		return
	}
	client := db.MongoCl.Client
	collection := client.Database("photos").Collection("user_photos")

	var curRec userPhoto
	filter := bson.M{"photo_url": res.PhotoURL}
	collection.FindOne(ctx, filter).Decode(&curRec)

	if curRec.Sickness != "" {
		ctx.JSON(200, gin.H{"result": gin.H{"photo_url": curRec.Photo_url, "sickness": curRec.Sickness, "accuracy": curRec.Accuracy}})
		return
	}

	ctx.JSON(200, gin.H{"result": gin.H{"sickness": "😸", "accuracy": "😸"}})
}

func PostResult(ctx *gin.Context) {
	var res Result
	if err := ctx.BindJSON(&res); err != nil {
		fmt.Println("error decoding result: " + err.Error())
		return
	}

	client := db.MongoCl.Client
	collection := client.Database("photos").Collection("user_photos")
	filter := bson.M{"photo_url": res.PhotoURL}

	var curRec userPhoto
	err := collection.FindOne(ctx, filter).Decode(&curRec)
	if err != nil {
		fmt.Println("error decoding userPhoto: ", err)
		return
	}

	curRec.Sickness = res.Sickness
	updateRec := bson.M{
		"$set": bson.M{
			"accuracy": res.Accuracy,
			"sickness": res.Sickness,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, updateRec)
	if err != nil {
		fmt.Println("error updating userPhoto record with results: ", err)
		return
	}

	fmt.Printf("the response is: %s with accuracy of %s\n", res.Sickness, res.Accuracy)
}
