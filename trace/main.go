package main

import (
	"fmt"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/pkg/profile"

	memoryProfiling "github.com/shadowpeak100/GoMemoryProfiling"
)

func main() {
	defer fmt.Println("Trace profile finished, access with: 'go tool trace trace.out'")
	defer profile.Start(profile.TraceProfile, profile.ProfilePath("profiling")).Stop()

	// Code you want to profile
	wg := new(sync.WaitGroup)
	for i := 0; i < memoryProfiling.Workers; i++ {
		wg.Add(1)
		go work(wg.Done)
	}
	wg.Wait()
}

func work(done func()) {
	// Simulate some work
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	done()
}
