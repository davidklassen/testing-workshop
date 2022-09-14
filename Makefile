.PHONY: init
init:
	./scripts/init.sh

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test:
	go test ./... -v

.PHONY: test-integration
test-integration:
	go clean -testcache && go test -p 1 ./tests/... -v -tags integration

.PHONY: test-all
test-all: test test-integration

.PHONY: start
start:
	docker-compose up -d --build --remove-orphans

.PHONY: stop
stop:
	docker-compose down

.PHONY: logs
logs:
	docker-compose logs -f
