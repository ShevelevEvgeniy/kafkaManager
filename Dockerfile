FROM golang:latest

WORKDIR /app

COPY ./ ./

RUN apt-get update && apt-get install -y \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

RUN curl -sL https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xz -C /usr/local/bin

RUN go mod download
RUN go build -o kafka_maneger ./cmd/kafka_maneger/main.go

CMD ["./kafka_maneger"]