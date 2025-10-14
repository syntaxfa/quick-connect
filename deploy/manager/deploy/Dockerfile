FROM golang:1.25.3 AS builder

ENV APP_NAME=manager

WORKDIR /$APP_NAME

# copy app files
COPY . .

# build application
RUN CGO_ENABLED=0 go build -o $APP_NAME cmd/manager/main.go

# run stage
FROM alpine:3.22.2 AS runtime

COPY --from=builder /manager/manager ./
COPY --from=builder /manager/deploy/manager/config.yml ./deploy/manager/config.yml
COPY --from=builder /manager/pkg/translation ./pkg/translation

CMD ["./manager", "start"]
