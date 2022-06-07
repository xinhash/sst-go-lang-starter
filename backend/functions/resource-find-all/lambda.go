package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Resource struct {
	PK string `dynamodbav:"PK" json:"category"`
	SK string `dynamodbav:"SK" json:"sub_category"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "ap-south-1"
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	tableName := os.Getenv("table")

	result, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("begins_with(PK, :pk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "Resource"},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItems: %s", err)
	}
	var resources []Resource
	err = attributevalue.UnmarshalListOfMaps(result.Items, &resources)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Records, %v", err))
	}
	body, err := json.Marshal(resources)
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Failed to unmarshal Records, %v", err))
	}
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
