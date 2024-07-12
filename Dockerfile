FROM golang

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY docker_config.json ./

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/memosproxy

EXPOSE 8080

HEALTHCHECK --interval=5m --timeout=3s CMD curl -f http://localhost:8080/ok || exit 1

CMD ["/app/memosproxy", "/app/docker_config.json"]
