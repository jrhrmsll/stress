# stress

stress is a little load generator program written in Go, with a "constant throughput" style and using a **HdrHistogram** to avoid the **Coordinate Omission** problem.

## Installation
````
docker build -t stress .
````

## Options
```
docker run -it stress -h
Usage of stress:
  -duration string
    	Test duration. (default "60s")
  -throughput int
    	Throughput (requests per second). (default 1)
  -timeout string
    	Request timeout. (default "30s")
  -url string
    	URL of the system under test. (default "http://localhost/")
```

e.g.
```
go run main.go -throughput 10 -duration 1s -u http://localhost:8080/bookmarks

        URL: http://localhost:8080/bookmarks
 Throughput: 10 request/s
   Interval: 100ms
    Timeout: 30s
   Requests: 10

 Response times percentiles:
       50th:      7ms
       75th:      9ms
       95th:     10ms
       99th:     10ms
     99.9th:     10ms
    99.95th:     10ms
    99.99th:     10ms

```

The program also create a CSV (results.csv) file. The first line contain the values corresponding to `throughput, duration, timeout, url, interval, requests`
followed by latencies values and their occurrences.

e.g.
```
10,1s,30s,http://localhost:8080/bookmarks,100,10
8,1
7,1
10,1
6,4
9,3
```

Those results can be analyzed with other tools or alternative HdrHistogram implementations.

## References
[1] https://bravenewgeek.com/everything-you-know-about-latency-is-wrong/
