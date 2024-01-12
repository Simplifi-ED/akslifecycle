// SPDX-License-Identifier: GPL-3.0
// Copyright Authors of Akslifecycle

package internal

import (
	"context"
	"os"

	"github.com/charmbracelet/log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

func LogIntoAzure() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Errorf("AZURE_SUBSCRIPTION_ID: %s", extractSubString(os.Getenv("AZURE_SUBSCRIPTION_ID")))
		log.Errorf("AZURE_TENANT_ID: %s", extractSubString(os.Getenv("AZURE_TENANT_ID")))
		log.Errorf("AZURE_CLIENT_ID: %s", extractSubString(os.Getenv("AZURE_CLIENT_ID")))
		log.Errorf("AZURE_CLIENT_SECRET: %s", extractSubString(os.Getenv("AZURE_CLIENT_SECRET")))
		log.Fatalf("Failed to create Azure credential: %v", err)

	}
	// Azure SDK Resource Management clients accept the credential as a parameter.
	// The client will authenticate with the credential as necessary.
	client, err := armsubscription.NewSubscriptionsClient(cred, nil)
	if err != nil {
		log.Fatalf("Failed to create Azure subscription client: %v", err)
	}
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	_, err = client.Get(context.TODO(), subscriptionID, nil)
	if err != nil {
		log.Fatalf("Failed to create Azure client: %v", err)
	} else {
		log.Info("Successfully logged into Azure and retrieved subscription.")
	}
}

func extractSubString(stringObj string) string {
	stringLength := len(stringObj) / 2
	return stringObj[:-stringLength]
}
