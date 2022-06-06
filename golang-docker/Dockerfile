FROM golang:alpine as builder
RUN apk add git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
#COPY resources/* ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk update && apk --no-cache add tzdata ca-certificates

WORKDIR /root/
RUN mkdir database

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/database/data.json database/
#COPY --from=builder /app/resources resources/

# Expose port 8080 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./main"]
