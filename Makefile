# Makefile

.PHONY: kill-port dev make-migration gen gh-wh-proxy dev-ui dev-agent migrate-down

id ?=

kill-port:
	@lsof -ti:8000 | xargs kill -9

make-migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost/siphon?sslmode=disable" down

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost/siphon?sslmode=disable" up

gh-wh-proxy:
	smee -u https://smee.io/ZdaCIAdCc7Z02P --port 8000 --path /api/github/hook

dev:
	$(MAKE) gh-wh-proxy & GOENV=local air -c ./.air.api.toml & wait

dev-agent:
	DOCKER_API_VERSION=1.45 air -c ./.air.agent.toml

dev-ui:
	cd web && pnpm dev

deploy-agent:
	docker buildx build --platform linux/amd64 --push -f Dockerfile.agent . -t sharith/ciphon-agent

gen:
	@protoc \
		--proto_path=proto "proto/job_runs.proto" \
		--go_out=internal/protogen --go_opt=paths=source_relative \
  		--go-grpc_out=internal/protogen --go-grpc_opt=paths=source_relative