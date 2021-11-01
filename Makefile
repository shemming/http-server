build:
	go build -o bin/main main.go

tidy: 
	go mod tidy

run:
	go run main.go

all: build tidy run