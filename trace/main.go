package main

import (
	"time"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.TraceProfile, profile.ProfilePath("profiling")).Stop()

	// Code you want to profile
	for i := 0; i < 1000; i++ {
		work()
	}
}

func work() {
	// Simulate some work
	time.Sleep(time.Millisecond * 10)
}
