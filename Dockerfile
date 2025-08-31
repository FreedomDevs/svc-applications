FROM golang:latest AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" .

FROM scratch
COPY --from=builder /app/svc-applications /svc-applications

EXPOSE 9003
ENV GIN_MODE=release

CMD ["/svc-applications"]
