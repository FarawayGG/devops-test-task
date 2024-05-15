ARG GOLANG_VERSION
FROM golang:${GOLANG_VERSION}-alpine3.19 AS builder
ENV CGO_ENABLED 0
ENV GOOS linux
COPY rootfs .
RUN cd src && go build -o /opt/main ./main.go

FROM alpine:3.19
USER 0
RUN adduser -D -u 10001 alpine
COPY --from=builder /opt/main /opt/main
WORKDIR /opt
EXPOSE 8080
USER 10001
CMD ["/opt/main"]
