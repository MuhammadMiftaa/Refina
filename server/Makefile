include .env

$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))

api:
	go run ./cmd/api/main.go

worker:
	go run ./cmd/worker/main.go

go:
	@trap 'kill 0' INT TERM EXIT; \
	go run ./cmd/api/main.go & \
	go run ./cmd/worker/main.go & \
	wait

migrate:
	@if [ -z "$(to)" ]; then \
		goose up; \
	else \
		goose up-to $(to); \
	fi

migration:
	@goose create $(name) sql

rollback:
	@if [ -z "$(to)" ]; then \
		goose down; \
	else \
		goose down-to $(to); \
	fi

migration-status:
	@goose status

seeder:
	@goose -dir ./config/db/seeder create $(name) sql

seed:
	@goose -dir ./config/db/seeder -no-versioning up

seed-reset:
	@goose -dir ./config/db/seeder -no-versioning reset