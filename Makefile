scrape:
	@go run ./cmd/scraper

server:
	@templ generate ./cmd/server && go run ./cmd/server