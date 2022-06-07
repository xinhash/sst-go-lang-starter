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

type ConfirmForgotPasswordInput struct {
	UserId           string `json:"user_id" validate:"required"`
	Email            string `json:"email" validate:"required"`
	ConfirmationCode string `json:"confirmation_code" validate:"required"`
	Password         string `json:"password"  validate:"required"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var passwordResetUser ConfirmForgotPasswordInput
	ApiResponse := events.APIGatewayProxyResponse{}
	err := json.Unmarshal([]byte(request.Body), &passwordResetUser)
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
		user := &cognito.ConfirmForgotPasswordInput{
			ClientId:         aws.String(os.Getenv("cognitoClientId")),
			ConfirmationCode: aws.String(passwordResetUser.ConfirmationCode),
			Username:         aws.String(passwordResetUser.Email),
			Password:         aws.String(passwordResetUser.Password),
		}

		_, err = cognitoClient.ConfirmForgotPassword(context.TODO(), user)

		if err != nil {
			body := "Error: Invalid JSON payload ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
			ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500, Headers: map[string]string{"Content-Type": "application/json"}}
			return ApiResponse, err
		}
		response := struct {
			PassowrdConfirmed bool `json:"passowrd_confirmed"`
		}{
			PassowrdConfirmed: true,
		}

		body, _ := json.Marshal(response)
		ApiResponse = events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200, Headers: map[string]string{"Content-Type": "application/json"}}
	}

	return ApiResponse, nil

}

func main() {
	lambda.Start(Handler)
}
