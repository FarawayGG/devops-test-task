FROM golang:1.21 as build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY ./files /
RUN go build -o ./app ./main.go

FROM alpine:3.15.0

WORKDIR /app

COPY --from=build /app/app /app/app

ENTRYPOINT ["./app"]
