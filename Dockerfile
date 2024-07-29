FROM golang:latest

WORKDIR /app

#COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
#RUN chmod +x /usr/local/bin/wait-for-it.sh

COPY ./ ./

RUN apt-get update && apt-get install -y \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

RUN curl -sL https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xz -C /usr/local/bin

RUN go mod download
RUN go build -o kafka_maneger ./cmd/kafka_maneger/main.go

#CMD ["wait-for-it.sh", "db:5432", "--", "./kafka_maneger"]
CMD ["./kafka_maneger"]