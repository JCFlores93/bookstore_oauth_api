FROM golang:latest as builder

LABEL maintaner="Jean Carlo Flores Carrasco <jeancarlo_flores93@hotmail.com>"

RUN mkdir /app

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/src

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch

COPY --from=builder /app/src/main /root/

WORKDIR /root

CMD ["./main"]