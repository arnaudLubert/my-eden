FROM golang:1.15 AS builder
LABEL build=golang-temp
WORKDIR /go
ENV GO111MODULE=on
ENV GOPATH=/
ADD src /go/src
COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o server ./src

FROM alpine:latest
RUN apk --no-cache add imagemagick mariadb-connector-c mariadb-client
WORKDIR /root/
COPY --from=builder /go/server .
RUN mkdir website_bridge /backups-sql
WORKDIR /root/website_bridge
CMD ["../server"]
