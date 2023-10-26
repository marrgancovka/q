FROM golang:1.21.1-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o ./bin ./cmd/main/main.go

EXPOSE 8080

CMD ["./bin"]

# docker build -t go-service .
# docker run -p 8080:8080  go-service