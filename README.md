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

## References
[1] https://bravenewgeek.com/everything-you-know-about-latency-is-wrong/
