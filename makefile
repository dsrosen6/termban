run:
	go run cmd/main/main.go

build:
	go build -o bin/termban cmd/main/main.go

log:
	tail -f termban.log