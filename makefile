lint:
	golangci-lint run ./internal/... ./cmd/... --build-tags="unit integration"

fmt:
	go fmt ./internal/... ./cmd/...

generate:
	go generate ./...

tests-unit: generate
	go test ./internal/... --tags=unit

tests-integration: generate
	./cmd/sh/tests-integration.sh

run-guideliner-debug:
	GUIDELINER_ENV=debug go run ./cmd/guildeliner/main.go

run-guideliner-development:
	GUIDELINER_ENV=development go run ./cmd/guildeliner/main.go

run-guideliner-production:
	GUIDELINER_ENV=production go run ./cmd/guildeliner/main.go

vendor:
	go mod vendor

build-guideliner: vendor
	go build -o ./bin/guideliner -mod vendor ./cmd/guideliner/main.go

build-migrations: vendor
	go build -o ./bin/migrations -mod vendor  ./cmd/migrations/main.go

build-clean-db: vendor
	go build -o ./bin/clean_db -mod vendor  ./cmd/clean_postgresql/main.go

build-all: vendor build-guideliner build-migrations build-clean-db

create-migration:
	touch ./internal/migrations/"$$(date +'%s_tempname.go')"

