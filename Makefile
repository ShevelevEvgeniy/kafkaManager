include .env

install:
	@$(MAKE) -s down
	@$(MAKE) -s docker-build
	@$(MAKE) -s up
	@$(MAKE) -s migrate-up
	@echo "--- Application installed ---"

up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans

docker-build:
	docker build -t kafka-manager .

shell:
	docker-compose exec kafka-manager /bin/bash -c "$(cmd)"

migrate-up:
	make shell cmd="migrate -source $(MIGRATION_URL) -database $(DB_DRIVER_NAME)://$(DB_USER_NAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE) -verbose up"

migrate-create:
	make shell cmd="migrate create -ext sql -dir migrations $(name)"

migrate-down:
	make shell cmd="migrate -source $(MIGRATION_URL) -database $(DB_DRIVER_NAME)://$(DB_USER_NAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE) -verbose down"

restart:
	@$(MAKE) -s docker-build
	@docker-compose up -d --no-deps --build kafka-manager

run-tests:
	go test -v ./...