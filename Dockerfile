# Build stage
FROM golang:1.23.2-alpine3.20 AS builder
LABEL intermediateStageToBeDeleted=true

WORKDIR /home/app

COPY /src/go.mod /src/go.sum ./
RUN go mod download

COPY /src ./

RUN go build -o betterreads ./cmd/main.go

# Test stage
FROM builder AS betterreads-test-stage

CMD ["go", "test", "-v", "./tests"]

# Run stage
FROM alpine:3.20

WORKDIR /home/app

COPY --from=builder /home/app/betterreads ./

ENTRYPOINT ["./betterreads"]
