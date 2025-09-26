SHELL = /bin/bash
scrape:
	@go build -o bin ./cmd/scraper

server:
	@templ generate ./cmd/server && CGO_ENABLED=0 GOOS=linux go build -o bin/ ./cmd/server

server-watch:
	@templ generate --watch --proxy="http://localhost:8080" --cmd="go run cmd/server/main.go ./server.config.yml"