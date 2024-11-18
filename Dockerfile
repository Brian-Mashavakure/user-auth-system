FROM golang:1.22.4
WORKDIR /bin

COPY go.mod ./

RUN go mod tidy

COPY . .

ENV PORT=8080


CMD ["go", "run", "./cmd/main"]