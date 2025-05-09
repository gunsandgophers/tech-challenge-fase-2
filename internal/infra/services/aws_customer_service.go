package services

import (
	"tech-challenge-fase-1/internal/core/entities"
	valueobjects "tech-challenge-fase-1/internal/core/value_objects"
)

type AwsCustomerService struct {
	serviceClient *CognitoClient
}

func NewAwsCustomerService(
	region string,
	userPoolId string,
) (*AwsCustomerService, error) {
	cognitoClient, err := NewCognitoClient(region, userPoolId)
	if err != nil {
		return nil, err
	}
	return &AwsCustomerService{
		serviceClient: cognitoClient,
	}, nil
}

func (a *AwsCustomerService) GetCustomerByCPF(cpf *valueobjects.CPF) (*entities.Customer, error) {
	user, err := a.serviceClient.GetUser(cpf.Value())
	if err != nil {
		return nil, err
	}
	return entities.RestoreCustomer(user.Id, user.Name, user.Email, user.Username)
}

func (a *AwsCustomerService) GetCustomerById(id string) (*entities.Customer, error) {
	user, err := a.serviceClient.GetUserBySub(id)
	if err != nil {
		return nil, err
	}
	return entities.RestoreCustomer(user.Id, user.Name, user.Email, user.Username)
}

func (a *AwsCustomerService) CreateCustomer(customer *entities.Customer) (*entities.Customer, error) {
	cognitoCreateUser := &CognitoCreateUser{
		Username: customer.GetCPF().Value(),
		Name: customer.GetName(),
		Email: customer.GetEmail().Value(),
	}
	user, err := a.serviceClient.CreateUser(cognitoCreateUser)
	if err != nil {
		return nil, err
	}
	return entities.RestoreCustomer(user.Id, user.Name, user.Email, user.Username)
}
