// SPDX-License-Identifier: GPL-3.0
// Copyright Authors of Akslifecycle

package lifecycle

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/Simplifi-ED/akslifecycle/internal"
	"github.com/charmbracelet/log"
)

func StartNode(clusterName *string, resourceGroup *string, nodepoolName *string) {
	internal.LogIntoAzure()

	log.Printf("Starting nodepool %s in cluster %s in resource group %s", *nodepoolName, *clusterName, *resourceGroup)
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAgentPoolsClient().BeginCreateOrUpdate(ctx, *resourceGroup, *clusterName, *nodepoolName, armcontainerservice.AgentPool{
		Properties: &armcontainerservice.ManagedClusterAgentPoolProfileProperties{
			PowerState: &armcontainerservice.PowerState{
				Code: to.Ptr(armcontainerservice.CodeRunning),
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	_ = res
}
