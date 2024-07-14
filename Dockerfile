# syntax=docker/dockerfile:1

FROM golang:1.22  AS builder

# Set destination for COPY
WORKDIR /app

COPY ./ /app

# Download Go modules
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /burgers-api /app/cmd/burgers-api/main.go

FROM scratch

COPY --from=builder /app /
COPY --from=builder /burgers-api /burgers-api

# Run
CMD [ "/burgers-api" ]
