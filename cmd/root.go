// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Akslifecycle

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
	"github.com/muandane/akslifecycle/internal"
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

var configFile string

var (
	config   Config
	cronJobs []*cron.Cron
	wg       sync.WaitGroup
)

var rootCmd = &cobra.Command{
	Use:   "akslifecycle",
	Short: "akslifecycle CLI",
	Long:  `akslifecycle is a cli tool to start & stop nodes with cron schedule`,
	Run: func(cmd *cobra.Command, args []string) {
		azureAuth := internal.NewAzureAuth()
		azureAuth.LogIntoAzure()

		viper.SetConfigFile(configFile)

		viper.WatchConfig()

		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Warnf("Config file got updated: %v", e.Name)

			err := viper.ReadInConfig()
			if err != nil {
				log.Fatalf("Failed to read config file: %v", err)
				return
			}

			err = viper.Unmarshal(&config)
			if err != nil {
				log.Errorf("Failed to decode into struct, %v", err)
				return
			}
			log.Info("Config is reloaded !")

			for _, c := range cronJobs {
				c.Stop()
			}
			cronJobs = nil

			wg.Add(len(config.Resources))
			for _, resource := range config.Resources {

				go func(resource Resource) {

					defer wg.Done()
					c := cron.New()

					clusterName := resource.ClusterName
					resourceGroup := resource.ResourceGroupName
					startSchedule := resource.StartSchedule
					stopSchedule := resource.StopSchedule

					// Define the schedule for starting the program
					_, err := c.AddFunc(startSchedule, func() {
						for _, nodepool := range resource.NodePools {
							nodepoolName := nodepool
							lifecycle.StartNode(&clusterName, &resourceGroup, &nodepoolName)
							log.Info("Waiting for next cron job...")
						}
					})

					if err != nil {
						log.Errorf("Failed to add cron job: %v", err)
					}
					// Define the schedule for stopping the program
					_, err = c.AddFunc(stopSchedule, func() {
						for _, nodepool := range resource.NodePools {
							nodepoolName := nodepool
							lifecycle.StopNode(&clusterName, &resourceGroup, &nodepoolName)
							log.Info("Waiting for next cron job...")
						}
					})
					if err != nil {
						log.Errorf("Failed to add cron job: %v", err)
					}

					// Start the cron scheduler
					c.Start()

					// Add the cron job to the slice
					cronJobs = append(cronJobs, c)
				}(resource)
			}

			// Wait for all goroutines to finish
			wg.Wait()
		})

		sigs := make(chan os.Signal, 1)
		// Register the channel to receive SIGINT signals
		signal.Notify(sigs, os.Interrupt)
		// Create a goroutine to handle the SIGINT signal
		go func() {
			<-sigs
			// Perform cleanup tasks here
			fmt.Println("Received SIGINT signal. Performing cleanup tasks...")
			// Exit the program
			os.Exit(1)
		}()
		wg.Add(1)
		wg.Wait()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file to run aks lifecycle")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
