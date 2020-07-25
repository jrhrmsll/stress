# stress

stress is a little load generator program written in Go, with a "constant throughput" style and using a **HdrHistogram** to avoid the **Coordinate Omission** problem.

## Installation
````
docker build -t stress .
````

## Run
```
docker run -it stress -t 10 -d 1s -u http://<ip>:<port>/
```

### Sample output

```
go run main.go -t 10 -d 1s -u http://localhost:8080/bookmarks

        URL: http://localhost:8080/bookmarks
 Throughput: 10 request/s
   Interval: 100ms
   Requests: 10

 Response times percentiles:
       50th:      7ms
       75th:      8ms
       95th:      8ms
       99th:      8ms
     99.9th:      8ms
    99.95th:      8ms
    99.99th:      8ms

```

## References
[1] https://bravenewgeek.com/everything-you-know-about-latency-is-wrong/
