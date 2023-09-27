FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the rest of your application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# Expose the port your application will run on
EXPOSE 8080

# Define the command to run your application
CMD ["/docker-gs-ping"]