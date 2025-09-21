FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main app/main.go

CMD ["sh", "./startup.sh"]
