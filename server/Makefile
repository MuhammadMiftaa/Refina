include .env

$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))

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