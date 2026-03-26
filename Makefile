.PHONY: build run clean docker-build docker-run

build:
	go build -o atp ./cmd/server

run: build
	./atp

clean:
	rm -f atp

docker-build:
	docker build -t azzurro-atp .

docker-run: docker-build
	docker-compose up