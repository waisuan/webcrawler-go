run:
	go run cmd/cli/main.go -targetUrl=$(targetUrl)

test:
	go test -v ./...