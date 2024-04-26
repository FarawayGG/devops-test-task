FROM golang:1.16-alpine AS builder
WORKDIR /opt/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.19.1
RUN apk --no-cache add ca-certificates
WORKDIR /opt/app/
COPY --from=builder /opt/app/main .
EXPOSE 8080
CMD ["./main"]
