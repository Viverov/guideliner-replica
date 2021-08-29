lint:
	golangci-lint run ./internal/... ./cmd/... --build-tags="unit integration"

fmt:
	go fmt ./internal/... ./cmd/...

generate:
	go generate ./...

tests-unit: generate
	go test ./internal/... --tags=unit

tests-integration: generate
	go test -p 1 ./internal/... --tags=integration

tests-integration-in-docker:
	./cmd/sh/run_integration_tests_by_docker.sh

guideliner-run-debug:
	GUIDELINER_ENV=debug go run ./cmd/guildeliner/main.go

guideliner-run-development:
	GUIDELINER_ENV=development go run ./cmd/guildeliner/main.go

guideliner-run-production:
	GUIDELINER_ENV=production go run ./cmd/guildeliner/main.go

vendor:
	go mod vendor

guideliner-build: vendor
	go build -o ./bin/guideliner -mod vendor ./cmd/guideliner/main.go

migrations-build: vendor
	go build -o ./bin/migrations -mod vendor  ./cmd/migrations/main.go

clean-db-build: vendor
	go build -o ./bin/clean_db -mod vendor  ./cmd/clean_postgresql/main.go

all-build: vendor guideliner-build migrations-build clean-db-build

create-migration:
	touch ./internal/migrations/"$$(date +'%s_tempname.go')"

