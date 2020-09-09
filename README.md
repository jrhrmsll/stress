# stress

stress is a little load generator program written in Go, with a "constant throughput" style and using a **HdrHistogram** to avoid the **Coordinate Omission** problem.

## Installation
````
go get github.com/jrhrmsll/stress
````

## Options
```
stress -h
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

If the [Bookmarks Service](https://github.com/jrhrmsll/bkd) is running, the following will
launch 10 request per second during 5s.

```
stress -throughput 10 -duration 5s -url http://localhost:8080/bookmarks

        URL: http://localhost:8080/bookmarks
 Throughput: 10 request/s
   Interval: 100ms
    Timeout: 30s
   Requests: 50

 Response times percentiles:
       50th:      7ms
       75th:      8ms
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
10,5s,30s,http://localhost:8080/bookmarks,100,50
6,9
9,1
5,3
4,1
10,3
7,23
8,10
```

These results can be analyzed with other tools or alternative HdrHistogram implementations.

## References
[1] https://bravenewgeek.com/everything-you-know-about-latency-is-wrong/
