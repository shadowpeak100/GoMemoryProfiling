package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"

	"github.com/pkg/profile"
)

const (
	workers   = 8
	speedDial = 10
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	// when this is enabled you can connect with the following:
	//go tool pprof -http :9090 http://<IP_Address>:7777/debug/pprof/heap
	//go tool pprof -http :9090 http://<IP_Address>:7777/debug/pprof/profile
	//go tool pprof -http :9090 http://<IP_Address>:7777/debug/pprof/block
	//go tool pprof -http :9090 http://<IP_Address>:7777/debug/pprof/goroutine
	//go tool pprof -http :9090 http://<IP_Address>:7777/debug/pprof/trace
	go func() { http.ListenAndServe(":7777", nil) }()

	//trace profiling
	//defer profile.Start(profile.TraceProfile, profile.ProfilePath("profiling")).Stop()
	//defer fmt.Println("Trace profile started, access after program has exited with: 'go tool trace trace.out'")

	//cpu profiling
	defer profile.Start(profile.CPUProfile, profile.ProfilePath("profiling")).Stop()
	defer fmt.Println("CPU profile started, access after program has exited with: 'go tool pprof -http=:8080 cpu.pprof'")

	//memory profiling
	//defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath("profiling")).Stop()
	//defer fmt.Println("Memory profile started, access after program has exited with: 'go tool pprof -http=:8080 mem.pprof'")

	fileChan := make(chan string)
	go loader(fileChan)

	output := make(chan string)
	var wg sync.WaitGroup

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			worker(fileChan, output)
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
			fmt.Println(i, "lines saved in strings builder")
		}
	}

	writeToFile("output/output.txt", builder.String())
}
