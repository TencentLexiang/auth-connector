FROM golang:alpine3.15 AS builder

WORKDIR /build

ENV GOPROXY https://goproxy.cn
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-connector main.go

FROM alpine:3.15 AS final

WORKDIR /app
COPY --from=builder /build/auth-connector /app/
EXPOSE 9000
ENTRYPOINT ["/app/auth-connector"]