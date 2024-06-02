FROM golang:1.22.1-alpine AS builder

WORKDIR /user/src/backend
COPY . .
COPY go.mod .
RUN go mod download
RUN go build -v -x -o /usr/local/bin/backend

FROM alpine:latest
COPY --from=builder /usr/local/bin/backend /usr/local/bin/backend
EXPOSE 8080
ENTRYPOINT [ "/usr/local/bin/backend" ]