package cmd

import (
	"fmt"
	"log"

	"os"

	"github.com/muandane/akslifecycle/internal"
	"github.com/muandane/akslifecycle/utils/lifecycle"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goji",
	Short: "Goji CLI",
	Long:  `Goji is a cli tool to generate conventional commits with emojis`,
	Run: func(cmd *cobra.Command, args []string) {
		c := cron.New()
		clusterName := os.Getenv("CLUSTER_NAME")
		resourceGroup := os.Getenv("RESOURCE_GROUP")
		nodepoolName := os.Getenv("NODEPOOL_NAME")
		startSchedule := os.Getenv("START_CLUSTER") // Every day at 8am
		stopSchedule := os.Getenv("STOP_CLUSTER")   // Every day at 5pm

		azureAuth := internal.NewAzureAuth()
		azureAuth.LogIntoAzure()

		_, err := c.AddFunc(startSchedule, func() {
			lifecycle.StartNode(&clusterName, &resourceGroup, &nodepoolName)
			log.Println("Waiting for next cron job...")
		})
		if err != nil {
			fmt.Println("Failed to add cron job:", err)
		}

		_, err = c.AddFunc(stopSchedule, func() {
			lifecycle.StopNode(&clusterName, &resourceGroup, &nodepoolName)
			log.Println("Waiting for next cron job...")
		})
		if err != nil {
			fmt.Println("Failed to add cron job:", err)
		}
		// Start the cron scheduler
		c.Start()

		// Keep the program running
		select {}
	},
}

func init() {
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
