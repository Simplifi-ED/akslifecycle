package internal

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

type AzureAuth struct {
	credential *azidentity.DefaultAzureCredential
	client     *armsubscription.SubscriptionsClient
}

func NewAzureAuth() *AzureAuth {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Failed to create Azure credential: %v", err)
	}

	client, err := armsubscription.NewSubscriptionsClient(cred, nil)
	if err != nil {
		log.Fatalf("Failed to create Azure subscription client: %v", err)
	}

	return &AzureAuth{
		credential: cred,
		client:     client,
	}
}

func (a *AzureAuth) LogIntoAzure() {
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if subscriptionID == "" {
		log.Fatalf("AZURE_SUBSCRIPTION_ID environment variable is not set")
	}

	_, err := a.client.Get(context.TODO(), subscriptionID, nil)
	if err != nil {
		log.Fatalf("Failed to log into Azure: %v", err)
	}

	fmt.Printf("Logged into Azure with subscription ID: %s\n", subscriptionID)
}
