FROM dockerhub.timeweb.cloud/golang:1.22 AS builder
WORKDIR /src
COPY ./src .
RUN go mod download -x && CGO_ENABLED=0 go build -o /bin/app .

FROM dockerhub.timeweb.cloud/alpine:3.20 AS final
COPY --from=builder /bin/app /bin/
EXPOSE 8080 8081
ENTRYPOINT [ "/bin/app" ]
