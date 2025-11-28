FROM golang:1.25.4 AS builder

ENV APP_NAME=all-in-one

WORKDIR /$APP_NAME

# copy app files
COPY . .

# build application
RUN CGO_ENABLED=0 go build -o $APP_NAME cmd/all-in-one/main.go

# run stage
FROM alpine:3.22.2 AS runtime

COPY --from=builder /all-in-one/all-in-one ./
COPY --from=builder /all-in-one/deploy/manager/config.yml ./deploy/manager/config.yml
COPY --from=builder /all-in-one/deploy/chat/config.yml ./deploy/chat/config.yml
COPY --from=builder /all-in-one/deploy/notification/config.yml ./deploy/notification/config.yml
COPY --from=builder /all-in-one/deploy/admin/config.yml ./deploy/admin/config.yml

CMD ["./chat", "start"]
