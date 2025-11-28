FROM golang:1.25.4 AS builder

ENV APP_NAME=chat

WORKDIR /$APP_NAME

# copy app files
COPY . .

# build application
RUN CGO_ENABLED=0 go build -o $APP_NAME cmd/chat/main.go

# run stage
FROM alpine:3.22.2 AS runtime

WORKDIR /app

COPY --from=builder /chat/chat ./
COPY --from=builder /chat/deploy/chat/config.yml ./deploy/chat/config.yml
COPY --from=builder /chat/app/chatapp/repository/migrations/*.sql ./app/chatapp/repository/migrations/
COPY --from=builder /chat/pkg/translation ./pkg/translation

CMD ["./chat", "start"]
