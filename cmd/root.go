package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/Simplifi-ED/akslifecycle/internal"
	"github.com/Simplifi-ED/akslifecycle/utils/lifecycle"
	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
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
var config Config
var cronScheduler = cron.New()
var wg sync.WaitGroup

var rootCmd = &cobra.Command{
	Use:   "akslifecycle",
	Short: "akslifecycle CLI",
	Long:  `akslifecycle is a cli tool to start & stop nodes with cron schedule`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.LogIntoAzure()

		setupSignalHandler()
		loadConfig()
		setupCronJobs()

		cronScheduler.Run()
	},
}

func init() {
	viper.AutomaticEnv() // Bind environment variables [0][1]
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file to run aks lifecycle")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfig() {
	viper.SetConfigFile(configFile)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Warnf("Config file got updated: %v", e.Name)
		reloadConfig()
	})
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Error reading config file: %v", err)
		return
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Errorf("Failed to decode into struct, %v", err)
		return
	}
}

func reloadConfig() {
	cronScheduler.Stop()
	wg.Wait()
	loadConfig()
	setupCronJobs()
	cronScheduler.Start()
	log.Warnf("Config file reloaded successfully.")
}

func setupCronJobs() {
	for _, resource := range config.Resources {
		resource := resource // Avoid capturing the loop variable [2]
		cronScheduler.AddFunc(resource.StartSchedule, func() {
			for _, nodepool := range resource.NodePools {
				lifecycle.StartNode(&resource.ClusterName, &resource.ResourceGroupName, &nodepool)
				log.Info("Waiting for next cron job...")
			}
		})
		cronScheduler.AddFunc(resource.StopSchedule, func() {
			for _, nodepool := range resource.NodePools {
				lifecycle.StopNode(&resource.ClusterName, &resource.ResourceGroupName, &nodepool)
				log.Info("Waiting for next cron job...")
			}
		})
	}
}

func setupSignalHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Println("Received SIGINT signal. Configuration will be reloaded.")
		reloadConfig()
	}()
}
