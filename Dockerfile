FROM golang:1.23 as build
WORKDIR /app
COPY . .
RUN go build -o bot ./cmd/bot
EXPOSE 4444
CMD ["/app/bot"]

# FROM alpine:latest
# COPY --from=build /app/bot /bot
# RUN chmod +x /bot
# COPY ./data /data
# EXPOSE 4444
# CMD ["/bot"]