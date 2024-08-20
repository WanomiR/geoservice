# geoservice
Case study project for GoKata Academy

## Metrics
### Number of requests
```bash
curl localhost:8888/metrics | grep 'total_number_of_requests'
```
### Request duration histogram
```bash
curl localhost:8888/metrics | grep 'api_request_duration_seconds'
```
### Cached data access duration
```bash
curl localhost:8888/metrics | grep 'cache_access_duration_seconds'
```
	
## Profiling
Ресурс url для получения данных по профилированию: `/<my-path>/pprof`

```
Count	Profile
9	allocs /mycustompath/pprof/allocs
0	block /mycustompath/pprof/block
0	cmdline /mycustompath/pprof/cmdline
4	goroutine /mycustompath/pprof/goroutine
9	heap /mycustompath/pprof/heap
0	mutex /mycustompath/pprof/mutex
0	profile /mycustompath/pprof/profile
10	threadcreate /mycustompath/pprof/threadcreate
0	trace /mycustompath/pprof/trace

full goroutine stack dump /mycustompath/debug/pprof/goroutine
```

## Performance Analysis
При активной нагрузке (в виде множества заппросов на один эндпонит в течение нескольких секунд) узких мест не обнаружено. Бóльшая часть процессорного времени пришлась на вызовы системы (66.6%), остальные процессы распределены равномерно.

## Profiling Command
To start profiling of http server run this command

```bash
go tool pprof "http://localhost:8888/debug/pprof/profile?seconds=30"
```
At terminal (pprof) is showing type help it will give the list of pprof commands as below

- `top`: It displays a list of the top functions that consume CPU time. It shows which functions are the most CPU-intensive.
- `list <function>`: list command followed by the name of a function to display the source code for that function, highlighting the lines where the most CPU time is spent.
- `web`: This command generates an interactive graphical visualization of the profile data. It opens a web browser with a graphical representation of the call graph, making it easier to identify performance bottlenecks.
- `web list`: This command combines the web and list commands. It generates a web-based visualisation and allows you to click on functions in the visualisation to see their source code with highlighted hotspots.
- `peek <function>`: The peek command displays a summary of the profile data for a specific function, including the percentage of total CPU time it consumes and the number of times it was called.
- `disasm <function>`: The `disasm` command displays the assembly code for a specific function. This can be useful for low-level performance analysis.
- `pdf: The pdf command generates a PDF file containing the call graph visualisation. This is useful for sharing profiling results or documenting performance improvements.
- `text`: The text command displays the profile data in text form, showing the top functions and their CPU usage. It's a simple textual representation of the profiling data.
- `topN <N>`: Use topN followed by a number (e.g., top10) to display the top N functions consuming CPU time. This can help you focus on the most significant bottlenecks.
- `raw`: The raw command displays the raw profiling data in a machine-readable format. This is useful for advanced analysis or automation.
- `tags`: The tags command displays all available pprof tags in the profile data. Tags can provide additional context for profiling results.
- `quit` or `exit`: Use either of these commands to exit the pprof interactive session.


Alternatively:

```bash
curl --output profile "http://localhost:8888/debug/pprof/profile?seconds=30"
go tool pprof -http localhost:3435 profile
```

