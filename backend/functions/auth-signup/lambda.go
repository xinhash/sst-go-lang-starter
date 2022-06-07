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
	// Email is the email decided by the user
	// at signup time. This field is not required but it could
	// be useful to have
	Email string `json:"email" validate:"required"`

	// Password is the password decided by the user
	// at signup time. This field is required and no signup
	// can work without this.
	// To create a secure password, contraints on this field are
	// it must contain an uppercase and lowercase letter,
	// a special symbol and a number.
	Password string `json:"password" validate:"required"`
}

type SignUpResponse struct {
	UserId    *string `json:"user_id"`
	Confirmed bool    `json:"confirmed"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var newUser User
	ApiResponse := events.APIGatewayProxyResponse{}
	err := json.Unmarshal([]byte(request.Body), &newUser)
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
		user := &cognito.SignUpInput{
			ClientId: aws.String(os.Getenv("cognitoClientId")),
			Username: aws.String(newUser.Email),
			Password: aws.String(newUser.Password),
		}

		result, err := cognitoClient.SignUp(context.TODO(), user)

		if err != nil {
			body := "Error: Invalid JSON payload ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
			ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500, Headers: map[string]string{"Content-Type": "application/json"}}
			return ApiResponse, err
		}

		response := &SignUpResponse{
			UserId:    result.UserSub,
			Confirmed: result.UserConfirmed,
		}

		body, _ := json.Marshal(response)
		ApiResponse = events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200, Headers: map[string]string{"Content-Type": "application/json"}}
	}

	return ApiResponse, nil

}

func main() {
	lambda.Start(Handler)
}
