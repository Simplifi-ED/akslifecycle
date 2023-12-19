package internal

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func NewAKSClient(subscriptionID string) (*containerservice.ManagedClustersClient, error) {
	// Get the AKS cluster credentials using Azure Active Directory authentication
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, err
	}

	client := containerservice.NewManagedClustersClient(subscriptionID)
	client.Authorizer = authorizer

	return &client, nil
}
