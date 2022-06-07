package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type ConfirmUserInput struct {
	UserId           string `json:"user_id" validate:"required"`
	Email            string `json:"email" validate:"required"`
	ConfirmationCode string `json:"confirmation_code" validate:"required"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var nonConfirmedUser ConfirmUserInput
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
		user := &cognito.ConfirmSignUpInput{
			ClientId:         aws.String(os.Getenv("cognitoClientId")),
			Username:         aws.String(nonConfirmedUser.Email),
			ConfirmationCode: aws.String(nonConfirmedUser.ConfirmationCode),
		}
		_, err = cognitoClient.ConfirmSignUp(context.TODO(), user)

		if err != nil {
			body := "Error: Invalid JSON payload ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
			ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500, Headers: map[string]string{"Content-Type": "application/json"}}
			return ApiResponse, err
		}
		response := struct {
			UserConfirmed bool `json:"user_confirmed"`
		}{
			UserConfirmed: true,
		}

		body, _ := json.Marshal(response)
		ApiResponse = events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200, Headers: map[string]string{"Content-Type": "application/json"}}
	}

	return ApiResponse, nil

}

func main() {
	lambda.Start(Handler)
}
