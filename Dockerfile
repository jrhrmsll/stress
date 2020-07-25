FROM       golang:alpine as build
RUN        apk add --no-cache ca-certificates
WORKDIR    /app
COPY       . .
RUN        go build -o stress main.go

FROM       alpine:latest
RUN        apk add --no-cache ca-certificates
COPY       --from=build /app/stress /usr/local/bin/stress
ENTRYPOINT [ "stress" ]