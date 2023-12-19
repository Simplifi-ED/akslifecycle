package cmd

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/muandane/akslifecycle/internal"
	"github.com/robfig/cron"
)

func lifeCycle() {
	c := cron.New()
	clusterName := os.Getenv("CLUSTER_NAME")
	resourceGroupName := os.Getenv("RESOURCE_GROUP")
	subscriptionID := os.Getenv("SUBID")
	startSchedule := os.Getenv("START_CLUSTER") // Every day at 8am
	stopSchedule := os.Getenv("STOP_CLUSTER")   // Every day at 5pm

	// Define the schedule for starting the program
	_, err := c.AddFunc(startSchedule, func() {
		internal.StartAKSCluster(&subscriptionID, &clusterName, &resourceGroupName) // Call the startAKSCluster function from the internal package
		log.Println("Waiting for next cron job...")
	})
	if err != nil {
		log.Fatal("Failed to add cron job:", err)
	}

	// Define the schedule for stopping the program
	_, err = c.AddFunc(stopSchedule, func() {
		internal.StopAKSCluster(&subscriptionID, &clusterName, &resourceGroupName) // Call the stopAKSCluster function from the internal package
		log.Println("Waiting for next cron job...")
	})
	if err != nil {
		log.Fatal("Failed to add cron job:", err)
	}

	// MEm profile
	f, err := os.Create("mem.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		panic(err)
	}

	// CPU PROFILE
	f, err = os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	// Start the cron scheduler
	c.Start()

	// Keep the program running
	select {}
}
