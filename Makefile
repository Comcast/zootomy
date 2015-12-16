.PHONY: all

all:
	go build -v zkcfg.go
	docker build --no-cache -t zookeeper:3.5.0-alpha .
