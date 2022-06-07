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

type User struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var userInput User
	ApiResponse := events.APIGatewayProxyResponse{}
	err := json.Unmarshal([]byte(request.Body), &userInput)
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
		user := &cognito.InitiateAuthInput{
			ClientId: aws.String(os.Getenv("cognitoClientId")),
			AuthFlow: "USER_PASSWORD_AUTH",
			AuthParameters: map[string]string{
				"USERNAME": userInput.Email,
				"PASSWORD": userInput.Password,
			},
			ClientMetadata: map[string]string{
				"UserPoolId": os.Getenv("cognitoUserPoolId"),
			},
		}

		result, err := cognitoClient.InitiateAuth(context.TODO(), user)

		if err != nil {
			body := "Error: Invalid JSON payload ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
			ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500, Headers: map[string]string{"Content-Type": "application/json"}}
			return ApiResponse, err
		}

		response := struct {
			AccessToken  *string `json:"access_token"`
			RefreshToken *string `json:"refresh_token"`
			TokenType    *string `json:"token_type"`
			ExpiresIn    int32   `json:"expires_in"`
			IdToken      *string `json:"id_token"`
		}{
			AccessToken:  result.AuthenticationResult.AccessToken,
			RefreshToken: result.AuthenticationResult.RefreshToken,
			TokenType:    result.AuthenticationResult.TokenType,
			IdToken:      result.AuthenticationResult.IdToken,
			ExpiresIn:    result.AuthenticationResult.ExpiresIn,
		}

		body, _ := json.Marshal(response)
		ApiResponse = events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200, Headers: map[string]string{"Content-Type": "application/json"}}
	}

	return ApiResponse, nil

}

func main() {
	lambda.Start(Handler)
}
