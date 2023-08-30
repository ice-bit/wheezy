FROM golang:latest
# Prepare working environment
RUN mkdir /app
WORKDIR /app
# Copy source files
COPY . .
# Build and test the application
RUN go build -o wheezy
RUN go test ./... -v
# Run the app
EXPOSE 9000
CMD ["./wheezy"]