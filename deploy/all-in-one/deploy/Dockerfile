FROM golang:1.25.5 AS builder

ENV APP_NAME=all-in-one

WORKDIR /$APP_NAME

# copy app files.
COPY . .

# build application
RUN CGO_ENABLED=0 go build -o $APP_NAME cmd/all-in-one/main.go

# run stage
FROM alpine:3.22.2 AS runtime

WORKDIR /app

COPY --from=builder /all-in-one/all-in-one ./

# Copy configs
COPY --from=builder /all-in-one/deploy/manager/config.yml ./deploy/manager/config.yml
COPY --from=builder /all-in-one/deploy/chat/config.yml ./deploy/chat/config.yml
COPY --from=builder /all-in-one/deploy/notification/config.yml ./deploy/notification/config.yml
COPY --from=builder /all-in-one/deploy/admin/config.yml ./deploy/admin/config.yml
COPY --from=builder /all-in-one/app/managerapp/repository/migrations/*.sql ./app/managerapp/repository/migrations/
COPY --from=builder /all-in-one/app/chatapp/repository/migrations/*.sql ./app/chatapp/repository/migrations/
COPY --from=builder /all-in-one/app/notificationapp/repository/migrations/*.sql ./app/notificationapp/repository/migrations/
COPY --from=builder /all-in-one/pkg/translation ./pkg/translation
COPY --from=builder /all-in-one/app/adminapp/templates/ ./app/adminapp/templates/
COPY --from=builder /all-in-one/app/adminapp/static/ ./app/adminapp/static/

CMD ["./all-in-one", "start"]
