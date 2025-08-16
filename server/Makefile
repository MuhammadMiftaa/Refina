include .env.local

$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))

api:
	APP_ENV_FILE=.env.local go run ./cmd/api/main.go

worker:
	APP_ENV_FILE=.env.local go run ./cmd/worker/main.go

go:
	@trap 'kill 0' INT TERM EXIT; \
	APP_ENV_FILE=.env.local go run ./cmd/api/main.go & \
	APP_ENV_FILE=.env.local go run ./cmd/worker/main.go & \
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
	@goose -dir ./internal/seeder/sql create $(name) sql

seed:
	@goose -dir ./internal/seeder/sql -no-versioning up

seed-reset:
	@goose -dir ./internal/seeder/sql -no-versioning reset