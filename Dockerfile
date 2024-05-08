FROM golang:1.21-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:1.21-alpine as builder
COPY --from=modules /go/pkg /go/infra
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app

FROM alpine
COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app
COPY .env .env
RUN apk --no-cache add ca-certificates

CMD ["/app"]
