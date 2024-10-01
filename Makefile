.PHONY: run
run:
	@go run cmd/*.go

.PHONY: build
build:
	@go build -o ./app cmd/*.go

.PHONY: docker-build
docker-build:
	@docker build -t garantex -f Dockerfile .

.PHONY: test
test:
	@go test -v -race -cover ./...

.PHONY: lint
lint:
	@golint ./...

.PHONY: run-db
run-db:
	@docker run \
		-d \
		-v `pwd`/migrations:/docker-entrypoint-initdb.d/ \
		--rm \
		-p 5432:5432 \
		--name db \
		-e POSTGRES_DB=backend \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		postgres:12
