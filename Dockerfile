FROM golang:1.23 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/bot /app/bot
COPY ./data /app/data
EXPOSE 4444
ENTRYPOINT ["/app/bot"]
