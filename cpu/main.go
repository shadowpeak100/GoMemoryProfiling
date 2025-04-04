package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"

	"github.com/pkg/profile"

	memoryProfiling "github.com/shadowpeak100/GoMemoryProfiling"
)

func main() {
	go func() { http.ListenAndServe(":7777", nil) }()

	//cpu profiling
	defer fmt.Println("CPU profile finished, access with: 'go tool pprof -http=:8080 cpu.pprof'")
	defer profile.Start(profile.CPUProfile, profile.ProfilePath("profiling")).Stop()

	fileChan := make(chan string)
	go memoryProfiling.Loader(fileChan)

	output := make(chan string)
	var wg sync.WaitGroup

	wg.Add(memoryProfiling.Workers)
	for i := 0; i < memoryProfiling.Workers; i++ {
		go func() {
			memoryProfiling.Worker(fileChan, output)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		fmt.Println("closing output channel")
		close(output)
	}()

	var builder strings.Builder

	i := 0
	for word := range output {
		i++
		builder.WriteString(" " + word)
		if i%100_000 == 0 {
			fmt.Println(i, "lines saved in strngs builder")
		}
	}

	memoryProfiling.WriteToFile("output/output.txt", builder.String())
}
