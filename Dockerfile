## Build Stage
FROM golang:1.21 AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./


RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.buildcommit=$(git rev-parse --short HEAD) \
              -X main.buildtime=$(date '+%Y-%m-%dT%H:%M:%S%Z:00') \
              -s -w" \ 
    -o /go/bin/app

## Run Stage
FROM gcr.io/distroless/base-debian10 
COPY --from=build /go/bin/app /app
EXPOSE 8080
USER nonroot:nonroot
CMD ["/app"]