compile:
	GOOS=linux GOARCH=amd64 go build -o bin/pxe-linux cmd/pxe/main.go
packr_compile:
	GOOS=linux GOARCH=amd64 packr build -o bin/pxe-linux cmd/pxe/main.go
test:
	go test ./...