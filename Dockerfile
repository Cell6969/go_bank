# Build Stage
FROM golang:1.24.2-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go 
    
# Run Stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

# Debug: listing semua file
RUN echo "=== FILES IN /app ===" && ls -la /app && echo "=== END FILE LIST ==="

RUN chmod a+x start.sh wait-for.sh

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]