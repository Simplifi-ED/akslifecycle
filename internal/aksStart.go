package internal

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func startAKSCluster(subscriptionID, resourceGroup, clusterName string) error {
	// Get the AKS cluster credentials using Azure Active Directory authentication
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	client := containerservice.NewManagedClustersClient(subscriptionID)
	client.Authorizer = authorizer

	_, err = client.Start(resourceGroup, clusterName)
	if err != nil {
		return err
	}

	fmt.Printf("AKS cluster %s in resource group %s has been started\n", clusterName, resourceGroup)
	return nil
}
