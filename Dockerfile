FROM golang:1.21

RUN apt-get update && apt-get install -y \
    libfreetype6-dev \
    libpng-dev \
    libjpeg-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o clock-bot main.go

CMD ["/app/clock-bot"]
