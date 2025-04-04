# GoMemoryProfiling

The following is a demonstration on using various profiling tools within go. There are 3 different main folders, each demonstrating a different type of profiling: CPU, memory and trace.

## Profiling tools can be accessed either live or after file generation
### Viewing live:
In order to view the file live, the following should be imported `_ "net/http/pprof"` this 
allows for a side-effect import. We will not be calling pprof directly in this program but we 
rely on its functionality to serve web pages. At the top of the main function declare the following: 

    go func() { http.ListenAndServe(":7777", nil) }()

As the program is running you can step in and view the file live! The file generated is *a snapshot* and
if new information is required, run the command to access the webpage again inside terminal. To access the 
live page, run one of the following in the terminal:

    go tool pprof -http :9090 http://<IP_Address_Here>:7777/debug/pprof/goroutine <optional: path to binary executable>
    go tool pprof -http :9090 http://<IP_Address_Here>:7777/debug/pprof/heap <optional: path to binary executable>
    go tool pprof -http :9090 http://<IP_Address_Here>:7777/debug/pprof/threadcreate <optional: path to binary executable>
    go tool pprof -http :9090 http://<IP_Address_Here>:7777/debug/pprof/block <optional: path to binary executable>
    go tool pprof -http :9090 http://<IP_Address_Here>:7777/debug/pprof/mutex <optional: path to binary executable>
    go tool pprof -http :9090 http://<IP_Address_Here>:7777/debug/pprof/profile <optional: path to binary executable>
    go tool pprof -http :9090 http://<IP_Address_Here>:7777/debug/pprof/trace?seconds=5 <optional: path to binary executable>

This will bring up a webpage where the snapshot statistics are available for review in various formats. An optional argument,
the source file will allow additional statistics to be reviewed such as the memory usage per line.

### Generating files:

In the main function of the program you want to analyze, start a profile for the respective service you want to generate
statistics on. In this function, we start the profiling and defer the stop. Although this is all one one line, the defer statement
only get applied to the final function, which is stop here. Here are some sample commands:

    defer profile.Start(profile.TraceProfile, profile.ProfilePath("profiling")).Stop()
	defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath("profiling")).Stop()
    defer profile.Start(profile.CPUProfile, profile.ProfilePath("profiling")).Stop()

Notice on the memory profiler that we have a rate set. The rate right now is set to a very low number, 
allowing more samples to be taken. This will catch more objects and allocations but will run slower. It is perfectly
acceptable to run this with a larger number like 4000 and see what results are yielded, then turning the number down 
if objects or allocations are not being found.


### Viewing generated file:

After the file has been generated, it can be viewed with the following commands:

    go tool pprof -http=:8080 cpu.pprof
    go tool pprof -http=:8080 mem.pprof
    go tool trace trace.out

The port number (here 8080) can be changed to whatever port is desired or not being used.
the last argument, the file name should be a path to where the file is located. 
Here we are assuming that the command is being run from the folder where the file is located

The go tool trace will be brought up using chromium tools, for best support ensure the application
opens in chrome or another chromium based browser


