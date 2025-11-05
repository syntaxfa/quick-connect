FROM golang:1.25.4 AS builder

ENV APP_NAME=chat

WORKDIR /$APP_NAME

# copy app files
COPY . .

# build application
RUN CGO_ENABLED=0 go build -o $APP_NAME cmd/chat/main.go

# run stage
FROM alpine:3.22.2 AS runtime

COPY --from=builder /chat/chat ./
COPY --from=builder /chat/deploy/chat/config.yml ./deploy/chat/config.yml

CMD ["./chat", "start"]
