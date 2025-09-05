build:
	go build -o bin/conti ./cmd/main.go
run: build
	./bin/conti