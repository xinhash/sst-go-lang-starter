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
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var nonConfirmedUser ForgotPasswordInput
	ApiResponse := events.APIGatewayProxyResponse{}
	err := json.Unmarshal([]byte(request.Body), &nonConfirmedUser)
	if err != nil {
		body := "Error: Invalid JSON payload ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
		ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500, Headers: map[string]string{"Content-Type": "application/json"}}
		return ApiResponse, err
	} else {
		cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
			o.Region = "ap-south-1"
			return nil
		})
		if err != nil {
			panic(err)
		}
		cognitoClient := cognito.NewFromConfig(cfg)
		// FIXME: Check if user already confirmed
		user := &cognito.ForgotPasswordInput{
			ClientId: aws.String(os.Getenv("cognitoClientId")),
			Username: aws.String(nonConfirmedUser.Email),
		}
		output, err := cognitoClient.ForgotPassword(context.TODO(), user)
		log.Printf("output %v", output)
		if err != nil {
			body := "Error: Invalid JSON payload ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
			ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500, Headers: map[string]string{"Content-Type": "application/json"}}
			return ApiResponse, err
		}
		response := struct {
			CodeDeliveryStatus string `json:"code_delivery_status,omitempty"`
		}{
			CodeDeliveryStatus: "Please check your " + string(output.CodeDeliveryDetails.DeliveryMedium) + " for reset instructions",
		}

		body, _ := json.Marshal(response)
		ApiResponse = events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200, Headers: map[string]string{"Content-Type": "application/json"}}
	}

	return ApiResponse, nil

}

func main() {
	lambda.Start(Handler)
}
