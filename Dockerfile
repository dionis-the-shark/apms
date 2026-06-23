FROM golang:1.22-alpine AS build

WORKDIR /src

RUN apk add --no-cache git

COPY apms-env/go.mod ./apms-env/
COPY apms-task-tracker/go.mod apms-task-tracker/go.sum ./apms-task-tracker/

WORKDIR /src/apms-task-tracker
RUN go mod download

WORKDIR /src
COPY apms-env ./apms-env
COPY apms-task-tracker ./apms-task-tracker

WORKDIR /src/apms-task-tracker
RUN go build -o /out/server ./cmd/server

FROM alpine:3.20

WORKDIR /app
COPY --from=build /out/server /app/server
COPY apms-task-tracker/.env /app/.env

EXPOSE 8080

CMD ["/app/server"]
