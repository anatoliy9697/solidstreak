run:
	cd solidstreak-frontend && npm install && npm run build
	cd solidstreak-backend && go run cmd/main.go