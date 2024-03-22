scrape:
	@go build -o bin ./cmd/scraper

server:
	@templ generate ./cmd/server && go build -o bin/ ./cmd/server

server-watch:
	@templ generate --watch --proxy="http://localhost:8080" --cmd="go run cmd/server/main.go"