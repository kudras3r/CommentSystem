
FROM golang:1.23.0-alpine as builder


WORKDIR /app


COPY ../go.mod ../go.sum ./
RUN go mod tidy


COPY ../ .
COPY ../.env /app/.env 
COPY ../internal/storage/migration/migrations /app/internal/storage/migration/migrations




WORKDIR /app/cmd


RUN go build -o main main.go


FROM alpine:latest


RUN apk --no-cache add ca-certificates


COPY --from=builder /app/cmd/main /app/cmd/main
COPY --from=builder /app/.env /app/.env 
COPY --from=builder /app/internal/storage/migration/migrations /app/internal/storage/migration/migrations


WORKDIR /app/cmd


EXPOSE 8080


CMD ["./main", "--storage=inmemory"]
