package cmd

import (
	"fmt"
	"os"

	"github.com/muandane/akslifecycle/utils/lifecycle"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Resource struct {
	ResourceGroupName string
	ClusterName       string
	NodePools         []string
	StartSchedule     string
	StopSchedule      string
}

type Config struct {
	Resources []Resource
}

var rootCmd = &cobra.Command{
	Use:   "akslifecycle",
	Short: "akslifecycle CLI",
	Long:  `akslifecycle is a cli tool to start & stop nodes with cron schedule`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println("Failed to read config file:", err)
		}
		var config Config
		err = viper.Unmarshal(&config)
		if err != nil {
			panic(fmt.Errorf("unable to decode into struct, %v", err))
		}

		for _, resource := range config.Resources {
			go func(resource Resource) {
				c := cron.New()
				clusterName := resource.ClusterName
				resourceGroup := resource.ResourceGroupName
				startSchedule := resource.StartSchedule
				stopSchedule := resource.StopSchedule

				// Define the schedule for starting the program
				_, err := c.AddFunc(startSchedule, func() {
					for _, nodepool := range resource.NodePools {
						lifecycle.StartNode(&clusterName, &resourceGroup, &nodepool)
						fmt.Println("Waiting for next cron job...")
					}
				})
				if err != nil {
					fmt.Println("Failed to add cron job:", err)
				}
				// Define the schedule for stopping the program
				_, err = c.AddFunc(stopSchedule, func() {
					for _, nodepool := range resource.NodePools {
						lifecycle.StopNode(&clusterName, &resourceGroup, &nodepool)
						fmt.Println("Waiting for next cron job...")
					}
				})
				if err != nil {
					fmt.Println("Failed to add cron job:", err)
				}

				// Start the cron scheduler
				c.Start()

			}(resource)
		}
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
