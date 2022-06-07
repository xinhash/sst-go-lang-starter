package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

	var params = make([]string, 2)
	params = strings.Split(request.PathParameters["id"], "#")
	PK := params[0]
	SK := params[1]
	key := Resource{PK: "RESOURCE#" + PK, SK: SK}

	avs, err := attributevalue.MarshalMap(key)
	if err != nil {
		panic(err)
	}

	result, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       avs,
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}
	resource := Resource{}
	err = attributevalue.UnmarshalMap(result.Item, &resource)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	body, err := json.Marshal(resource)
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
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
