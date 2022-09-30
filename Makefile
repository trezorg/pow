start:
	docker-compose up -d

stop:
	docker-compose down --remove-orphans  -v

logs:
	docker-compose logs -f

test:
	go test -v ./...

build:
	go build -a -o pow ./cmd

.PHONY: test start stop build
