FROM golang:1.22 as builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /build/ratelimiter cmd/main.go

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/ratelimiter .
COPY --from=builder /build/cmd/.env .

EXPOSE 8080

ENTRYPOINT ["/app/ratelimiter"]