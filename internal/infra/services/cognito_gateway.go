package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type CognitoClient struct {
	client *cognito.CognitoIdentityProvider
	userPoolId string
}

type CognitoUser struct {
	Id string
	Username string
	Name string
	Email string
}

type CognitoCreateUser struct {
	Username string
	Name string
	Email string
}

func NewCognitoClient(region string, userPoolId string) (*CognitoClient, error) {
	conf := &aws.Config{
		Region: aws.String(region),
	}
	sess, err := session.NewSession(conf)
	if err != nil {
		return nil, err
	}
	client := cognito.New(sess)
	return &CognitoClient{
		client: client,
		userPoolId: userPoolId,
	}, nil
}

func (cc *CognitoClient) GetUserBySub(sub string) (*CognitoUser, error) {
	input := &cognito.ListUsersInput{
		UserPoolId: aws.String(cc.userPoolId),
		Filter: aws.String(fmt.Sprintf("sub = \"%s\"", sub)),
		Limit: aws.Int64(1),
	}
	output, err := cc.client.ListUsers(input)
	if err != nil {
		return nil, err
	}

	if len(output.Users) == 0 {
		return nil, fmt.Errorf("User not found")
	}

	user := &CognitoUser{
		Username: *output.Users[0].Username,
	}

	for _, attr := range output.Users[0].Attributes {
		if *attr.Name == "email" {
			user.Email = *attr.Value
		} else if *attr.Name == "sub" {
			user.Id = *attr.Value
		} else if *attr.Name == "name" {
			user.Name = *attr.Value
		}
	}
	return user, nil
}

func (cc *CognitoClient) GetUser(username string) (*CognitoUser, error) {
	input := &cognito.AdminGetUserInput{
		UserPoolId: aws.String(cc.userPoolId),
		Username: aws.String(username),
	}
	output, err := cc.client.AdminGetUser(input)
	if err != nil {
		return nil, err
	}

	user := &CognitoUser{
		Username: *output.Username,
	}

	for _, attr := range output.UserAttributes {
		if *attr.Name == "email" {
			user.Email = *attr.Value
		} else if *attr.Name == "sub" {
			user.Id = *attr.Value
		} else if *attr.Name == "name" {
			user.Name = *attr.Value
		}
	}
	return user, nil
}

func (cc *CognitoClient) CreateUser(user *CognitoCreateUser) (*CognitoUser, error) {

	input := &cognito.AdminCreateUserInput{
		MessageAction: aws.String("SUPPRESS"),
		UserAttributes: []*cognito.AttributeType{
			{Name: aws.String("email"), Value: aws.String(user.Email)},
			{Name: aws.String("name"), Value: aws.String(user.Name)},
		},
		Username: &user.Username,
		UserPoolId: aws.String(cc.userPoolId),
	}
	output, err := cc.client.AdminCreateUser(input)
	if err != nil {
		return nil, err
	}
	cognitoUser := &CognitoUser{
		Username: *output.User.Username,
	}
	for _, attr := range output.User.Attributes {
		if *attr.Name == "email" {
			cognitoUser.Email = *attr.Value
		} else if *attr.Name == "sub" {
			cognitoUser.Id = *attr.Value
		} else if *attr.Name == "name" {
			cognitoUser.Name = *attr.Value
		}
	}
	return cognitoUser, nil
}

