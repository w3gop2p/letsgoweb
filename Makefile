build:
	@go build -o bin/letsgoweb ./cmd/web
run: build
	@./bin/letsgoweb
test :
	@go test -v ./...
