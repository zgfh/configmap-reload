FROM golang:1.14.3-alpine3.11 as builder
ENV GO111MODULE on
ENV GOPROXY https://goproxy.io
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GO15VENDOREXPERIMENT 1

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o /bin/app *.go

FROM busybox

USER 65534

ARG BINARY=configmap-reload
COPY --from=builder /bin/app  /bin/configmap-reload

ENTRYPOINT ["/configmap-reload"]
