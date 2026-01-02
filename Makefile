.PHONY: songsee

songsee:
	mkdir -p bin
	go build -o bin/songsee ./cmd/songsee
