FROM golang:1.26.3

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o server .

EXPOSE 8080

CMD ["./server"]