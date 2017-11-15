get :
	go get ./...

test:
	go test -v ./...

run: test
	go run main.go

build: test
	go build