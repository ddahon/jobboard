FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go install github.com/a-h/templ/cmd/templ@v0.2.543
RUN make server && make scraper

FROM scratch
WORKDIR /app
COPY --from=builder /app/bin/server /app/server
COPY ./server.config.yml /app/config.yml
ENTRYPOINT ["/app/server", "/app/config.yml"]
