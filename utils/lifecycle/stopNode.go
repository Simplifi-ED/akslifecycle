package lifecycle

import (
	"fmt"
	"log"
	"os/exec"
)

func StopNode(clusterName *string, resourceGroup *string, nodepoolName *string) {
	log.Printf("Stopping Cron for nodepool %s in cluster %s in resource group %s\n", *nodepoolName, *clusterName, *resourceGroup)
	stopCmd := exec.Command("az", "aks", "nodepool", "stop", "--resource-group", *resourceGroup, "--cluster-name", *clusterName, "--name", *nodepoolName)
	output, err := stopCmd.CombinedOutput()
	if err != nil {
		log.Println("Failed to stop AKS nodepool:", err)
	}
	fmt.Println("Command output:", string(output))
}
