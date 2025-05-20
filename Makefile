build:
	@go build -o bin/ticketbot cmd/bot/main.go

run: build
	@bin/ticketbot