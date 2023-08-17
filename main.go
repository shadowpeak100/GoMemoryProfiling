package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"strings"
	"sync"

	"github.com/pkg/profile"
)

const (
	workers = 8
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	//go func() { http.ListenAndServe(":7777", nil) }()
	//go func() { http.ListenAndServe(":7778", nil) }()
	//go func() { http.ListenAndServe(":7779", nil) }()
	defer profile.Start(profile.TraceProfile, profile.ProfilePath("profiling")).Stop()
	fmt.Println("Trace profile started, access after program has exited with: 'go tool trace trace.out'")

	//amdahls law
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath("profiling")).Stop()
	//fmt.Println("CPU profile started, access after program has exited with: 'go tool pprof -http=:8080 cpu.prof'")

	//defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath(".")).Stop()
	//fmt.Println("Memory profile started, access after program has exited with: 'go tool pprof -http=:8080 mem.pprof'")

	//if *cpuprofile != "" {
	//	f, err := os.Create(*cpuprofile)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	pprof.StartCPUProfile(f)
	//	defer pprof.StopCPUProfile()
	//}

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
