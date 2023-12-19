package lifecycle

import (
	"fmt"
	"log"
	"os/exec"
)

func StartNode(clusterName *string, resourceGroup *string, nodepoolName *string) {
	log.Printf("Starting Cron for nodepool %s in cluster %s in resource group %s\n", *nodepoolName, *clusterName, *resourceGroup)
	startCmd := exec.Command("az", "aks", "nodepool", "start", "--resource-group", *resourceGroup, "--cluster-name", *clusterName, "--name", *nodepoolName)
	output, err := startCmd.CombinedOutput()
	if err != nil {
		log.Println("Failed to start AKS nodepool:", err)
	}
	fmt.Println("Command output:", string(output))
}
