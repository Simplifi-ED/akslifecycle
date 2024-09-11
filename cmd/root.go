package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Simplifi-ED/akslifecycle/utils/lifecycle"
	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var location *time.Location

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

		setupSignalHandler()
		loadConfig()
		setupCronJobs()
		logNextOperations()
		cronScheduler.Run()
	},
}

func init() {
	viper.AutomaticEnv() // Bind environment variables [0][1]
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file to run aks lifecycle")
	tz := os.Getenv("TZ")
	var err error
	if tz != "" {
		location, err = time.LoadLocation(tz)
		if err != nil {
			log.Warnf("Invalid TZ environment variable. Defaulting to UTC: %v", err)
			location = time.UTC
		}
	} else {
		location = time.UTC
	}
	log.Infof("Using time zone: %s", location)
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
	log.Infof("Config file loaded successfully: %v", viper.ConfigFileUsed())
	if err := viper.Unmarshal(&config); err != nil {
		log.Errorf("Failed to decode into struct, %v", err)
		return
	}
	log.Infof("Config unmarshaled successfully")
}

func reloadConfig() {
	// Remove all existing cron jobs before stopping the scheduler
	for _, resource := range cronScheduler.Entries() {
		cronScheduler.Remove(resource.ID)
	}
	cronScheduler.Stop()
	wg.Wait()
	setupCronJobs()
	logNextOperations()
	cronScheduler.Start()
	log.Warnf("Config file reloaded successfully.")
}

func setupCronJobs() {
	for _, resource := range config.Resources {
		resource := resource // Avoid capturing the loop variable [2]
		addSchedule(resource.StartSchedule, func() { startNode(resource) })
		addSchedule(resource.StopSchedule, func() { stopNode(resource) })
	}
}

func addSchedule(spec string, cmd func()) {
	_, err := cronScheduler.AddFunc(spec, cmd)
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}
}

func startNode(resource Resource) {
	for _, nodepool := range resource.NodePools {
		lifecycle.StartNode(&resource.ClusterName, &resource.ResourceGroupName, &nodepool)
		log.Info("Waiting for next cron job...")
	}
	logNextOperations()
}

func stopNode(resource Resource) {
	for _, nodepool := range resource.NodePools {
		lifecycle.StopNode(&resource.ClusterName, &resource.ResourceGroupName, &nodepool)
		log.Info("Waiting for next cron job...")
	}
	logNextOperations()
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

func logNextOperations() {
	now := time.Now().In(location)
	for _, resource := range config.Resources {
		startSchedule, err := cron.ParseStandard(resource.StartSchedule)
		if err != nil {
			log.Errorf("Failed to parse start schedule for %s: %v", resource.ClusterName, err)
			continue
		}
		stopSchedule, err := cron.ParseStandard(resource.StopSchedule)
		if err != nil {
			log.Errorf("Failed to parse stop schedule for %s: %v", resource.ClusterName, err)
			continue
		}

		nextStart := startSchedule.Next(now)
		nextStop := stopSchedule.Next(now)

		log.Infof("Next operations for %s:", resource.ClusterName)
		log.Infof("  Next start: %s", nextStart.In(location).Format(time.RFC3339))
		log.Infof("  Next stop: %s", nextStop.In(location).Format(time.RFC3339))
	}
}
