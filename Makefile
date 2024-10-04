# Makefile

.PHONY: kill-port dev make-migration gh-wh-proxy dev-ui

kill-port:
	@lsof -ti:8000 | xargs kill -9

make-migration:
	migrate create -ext sql -dir migrations -seq $(name)

gh-wh-proxy:
	smee -u https://smee.io/ZdaCIAdCc7Z02P --port 8000 --path /api/github/hook

dev:
	$(MAKE) gh-wh-proxy & air & wait

dev-ui:
	cd web && pnpm dev