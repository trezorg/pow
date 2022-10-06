start:
	docker-compose up -d --build

stop:
	docker-compose down --remove-orphans  -v

logs:
	docker-compose logs -f

test:
	go test -v -count 1 ./...

build:
	go build -a -o pow ./cmd

.PHONY: test start stop build
