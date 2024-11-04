FROM golang:1.23 as build
WORKDIR /app
COPY . .
RUN go build -o /bot ./cmd/bot

FROM scratch
COPY --from=build /bot /bot
COPY --from=build ./data /data
EXPOSE 4444
CMD ["/bot"]