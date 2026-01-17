run:
	cd solidstreak-frontend && npm install && npm run build
	go mod tidy
	cd solidstreak-backend && go run cmd/main.go