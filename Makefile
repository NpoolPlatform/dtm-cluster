test:
	go test ./...

deps:
	go get -u github.com/ugorji/go/codec@latest
	go mod tidy -compat=1.17
