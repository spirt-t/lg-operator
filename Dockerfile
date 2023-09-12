ARG GOVERSION=1.19
FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV GOSUMDB=off
ENV GO111MODULE=on
ENV GOPRIVATE=""
ENV GONOPROXY=""
ENV CGO_ENABLED 0

ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

COPY config.yaml ./
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o lg-operator cmd/lg-operator/main.go

FROM alpine

COPY --from=builder /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
ENV TZ Europe/Moscow

WORKDIR /app
COPY --from=builder /build/lg-operator /app/lg-operator
COPY --from=builder /build/config.yaml /app/config.yaml
EXPOSE 7000
EXPOSE 7002
CMD ["/app/lg-operator"]
