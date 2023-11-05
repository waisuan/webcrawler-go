run:
	go run cmd/cli/main.go -targetUrl=$(targetUrl)

dev:
	APP_ENV="dev" go run cmd/cli/main.go -targetUrl=$(targetUrl)

test:
	APP_ENV="test" go test -v ./...