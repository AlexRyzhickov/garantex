.PHONY: run
run:
	@go run cmd/*.go -pg-user=postgres -pg-pass=postgres -pg-db=backend -pg-host=localhost -pg-port=5432

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

.PHONY: health-probe
health-probe:
	@grpc-health-probe -addr=:9090 -service garantex-proxy

.PHONY: gen-mocks
gen-mocks:
	@docker run -v `pwd`:/src -w /src vektra/mockery:v2.46.1 --case snake --dir internal --output internal/mock --outpkg mock --all --exported

.PHONY: load-test
load-test:
	@ab -c 4 -n 250 -k http://127.0.0.1:8080/price

