// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Akslifecycle

package lifecycle

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/log"
)

func StopNode(clusterName *string, resourceGroup *string, nodepoolName *string) {
	log.Printf("Stopping nodepool %s in cluster %s in resource group %s", *nodepoolName, *clusterName, *resourceGroup)
	stopCmd := exec.Command("az", "aks", "nodepool", "stop", "--resource-group", *resourceGroup, "--cluster-name", *clusterName, "--name", *nodepoolName)
	_, err := stopCmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to stop AKS nodepool: %v", err)
	}
	fmt.Printf("ResourceGroup: %s\n", *resourceGroup)
	fmt.Printf("Cluster: %s\n", *clusterName)
	fmt.Printf("Name: %s\n", *nodepoolName)
}
