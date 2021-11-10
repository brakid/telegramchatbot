FROM golang:latest AS builder

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

COPY . .

RUN go build -ldflags="-s -w" -o bot.out .

FROM alpine:latest

ENV BOT_API_KEY=setmeup
ENV CHAT_ID=12345
ENV MYSQL_URL=root:test@/spendings

WORKDIR /app

COPY --from=builder /app/bot.cfg .
COPY --from=builder /app/bot.out .

CMD ["./bot.out"]