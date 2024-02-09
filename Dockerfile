# dockerizing the backend API

# build stage
FROM golang:1.16-buster AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o goserver .


# Final stage
FROM debian:buster-slim
ENV PORT=8080
COPY --from=builder /app/goserver /bin/goserver
CMD ["/bin/goserver"]