FROM golang:1.24.5 AS builder

ENV APP_NAME=chat

WORKDIR /$APP_NAME

# copy app files
COPY . .

# build application
RUN CGO_ENABLED=0 go build -o $APP_NAME cmd/chat/main.go

# run stage
FROM alpine:3.22.1 AS runtime

COPY --from=builder /chat/chat ./
COPY --from=builder /chat/deploy/chat/config.yml ./deploy/chat/config.yml

CMD ["./chat", "start"]
