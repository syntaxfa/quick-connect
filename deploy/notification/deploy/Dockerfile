FROM golang:1.24.4 AS builder

ENV APP_NAME=notification

WORKDIR /$APP_NAME

# copy app files
COPY . .

# build application
RUN CGO_ENABLED=0 go build -o $APP_NAME cmd/notification/main.go

# run stage
FROM alpine:3.22.0 AS runtime

COPY --from=builder /notification/notification ./
COPY --from=builder /notification/deploy/notification/config.yml ./deploy/notification/config.yml
COPY --from=builder /notification/pkg/translation ./pkg/translation

CMD ["./manager", "start"]
