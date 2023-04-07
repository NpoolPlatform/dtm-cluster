test:
	go test ./...

deps:
	go mod tidy -compat=1.17
