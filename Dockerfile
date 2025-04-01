# Specifies a parent image
FROM golang:1.23 AS build

# Creates an app directory to hold your app’s source code
WORKDIR /app

COPY go.mod .
COPY go.sum .

# Installs Go dependencies
RUN go mod download

COPY . .

# Builds your app with optional configuration
RUN go build -o /bumflix-api ./cmd/main.go

FROM alpine:3.21 as run

# Copy the application executable from the build image
COPY --from=build /bumflix-api /bumflix-api

WORKDIR /app

EXPOSE 8080

CMD ["/bumflix-api"]
