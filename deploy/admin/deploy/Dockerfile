FROM golang:1.25.4 AS builder

ENV APP_NAME=admin

WORKDIR /$APP_NAME

# copy app files
COPY . .

# build application
RUN CGO_ENABLED=0 go build -o $APP_NAME cmd/admin/main.go

# run stage
FROM alpine:3.22.2 AS runtime

WORKDIR /app

COPY --from=builder /admin/admin ./
COPY --from=builder /admin/deploy/admin/config.yml ./deploy/admin/config.yml
COPY --from=builder /admin/pkg/translation ./pkg/translation

CMD ["./admin", "start"]
