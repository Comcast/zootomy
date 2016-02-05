.PHONY: all clean

all:
	go clean
	GOARCH=amd64 GOOS=linux go build -v zkcfg.go
	docker build --no-cache -t zookeeper:3.5.0-alpha .

clean:
	go clean
	docker rmi -f zookeeper:3.5.0-alpha
