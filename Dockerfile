# Build stage
FROM golang:1.23.0-alpine3.20 AS builder

WORKDIR /home/app

COPY /server/go.mod /server/go.sum ./
RUN go mod download

COPY /server ./



RUN go build -o service ./main.go

# Test stage
FROM builder AS test-stage

# CMD ["go", "test", "-cover", "-coverprofile=coverage/coverage.out", "./..."]
CMD ["sh", "-c", "go test -cover -coverprofile=coverage/coverage.out $(go list ./... | grep -v 'src/repository' | grep -v 'src/user-connector')"]


# Run stage
FROM alpine:3.20

WORKDIR /home/app

COPY --from=builder /home/app/service ./

ENTRYPOINT ["./service"]
