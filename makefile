.PHONY: setup-rigel
setup-rigel:
	cd script; ./config_rigel.sh
	
GOOSE_DIR = ./db/migrations/schema
GOOSE_DRIVER = sqlserver
GOOSE_DSN = "sqlserver://sa:SQLServer123@localhost:1433?database=CVL_KRA"

migrate-up:
	goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) $(GOOSE_DSN) up

migrate-down:
	goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) $(GOOSE_DSN) down

migrate-status:
	goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) $(GOOSE_DSN) status

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	goose -dir $(GOOSE_DIR) create $(name) sql


sqlc-generate:
	cd db; sqlc generate