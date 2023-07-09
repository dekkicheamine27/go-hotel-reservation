build:
    @go build -o bin/api
run:
    @./bin/api
test:
    @go test -v ./...