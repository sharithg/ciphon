# Makefile

.PHONY: kill-port dev make-migration gh-wh-proxy dev-ui migrate-down

id ?=

kill-port:
	@lsof -ti:8000 | xargs kill -9

make-migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost/siphon?sslmode=disable" down $(id)

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost/siphon?sslmode=disable" up

gh-wh-proxy:
	smee -u https://smee.io/ZdaCIAdCc7Z02P --port 8000 --path /api/github/hook

dev:
	$(MAKE) gh-wh-proxy & air & wait

dev-ui:
	cd web && pnpm dev