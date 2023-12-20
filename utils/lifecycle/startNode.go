// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Akslifecycle

package lifecycle

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/log"
)

func StartNode(clusterName *string, resourceGroup *string, nodepoolName *string) {
	log.Printf("Starting nodepool %s in cluster %s in resource group %s", *nodepoolName, *clusterName, *resourceGroup)
	startCmd := exec.Command("az", "aks", "nodepool", "start", "--resource-group", *resourceGroup, "--cluster-name", *clusterName, "--name", *nodepoolName)
	_, err := startCmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to start AKS nodepool: %v", err)
	}
	fmt.Printf("ResourceGroup: %s\n", *resourceGroup)
	fmt.Printf("Cluster: %s\n", *clusterName)
	fmt.Printf("Name: %s\n", *nodepoolName)
}
