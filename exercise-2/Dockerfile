FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

ENV DATABASE_URI="mongodb://mongodb:27017"

RUN go build -o /app/main ./cmd/main.go

CMD ["/app/main"]