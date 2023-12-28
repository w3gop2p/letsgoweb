build:
	@go build -o bin/letsgoweb
run: build
	@./bin/letsgoweb
test :
	@go test -v ./...
