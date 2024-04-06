FROM node:20-alpine as build-node

RUN mkdir -p /app/node_modules && chown -R node:node /app

WORKDIR /app

COPY web/package*.json ./

USER node

RUN npm install

COPY --chown=node:node web/. .

RUN npm run build

# Create a stage for building the application.
FROM golang:1.21 AS build-go

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

COPY --from=build-node /app/dist /app/internal/pkg/dist

RUN go env GOPATH | templ generate \
    && \
    CGO_ENABLED=0 GOOS=linux GARCH=amd64 go build -o /app/server cmd/chat-app.go


# Create a stage for the running image.
FROM alpine:latest AS final

WORKDIR /app

RUN apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

COPY --from=build-go /app/server /app

EXPOSE $PORT

ENTRYPOINT [ "/app/server -debug=false" ]
