scrape:
	@go build -o bin ./cmd/scraper

server:
	@templ generate ./cmd/server && go build -o bin/ ./cmd/server