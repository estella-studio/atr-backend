FROM golang:latest

WORKDIR /app

RUN useradd -ms /bin/bash app && chown -R app:app /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main app/main.go

CMD ["sh", "./startup.sh"]
