FROM golang:1.16-alpine

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o /app/build/service2

CMD ["/app/build/service2"]