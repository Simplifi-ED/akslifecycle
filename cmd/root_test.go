package cmd

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

func TestAddSchedule(t *testing.T) {
	// Setup
	// Create a new scheduler with second precision to make test fast
	cronScheduler = cron.New(cron.WithSeconds())
	executed := false
	spec := "* * * * * ?" // Every second

	// Define a test function that sets 'executed' to true
	testFunc := func() { executed = true }
	// Add the test function to the scheduler
	addSchedule(spec, testFunc)

	// Start the scheduler in a new goroutine
	go cronScheduler.Start()

	// Wait for the function to execute
	time.Sleep(2 * time.Second)

	// Stop the scheduler
	cronScheduler.Stop()

	// Check if the function was executed
	assert.True(t, executed, "Expected the scheduled function to be executed")
}
