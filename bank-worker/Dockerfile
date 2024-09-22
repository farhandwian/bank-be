# Use an official Golang runtime as the base image  
FROM golang:1.23-alpine3.20 AS BuildStage
  
# Set the working directory inside the container  
WORKDIR /go/src/bank-worker

COPY . .

# Build the Go app
RUN go build ./main.go


# Deploy Stage
FROM alpine:latest

COPY config/app.yml /

# Copy the built executable from BuildStage to the root directory
COPY --from=BuildStage /go/src/bank-worker /

ENTRYPOINT ["./main", "serve"]
