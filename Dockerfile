FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

EXPOSE 3000

CMD ["go", "run", "src/main.go"]
