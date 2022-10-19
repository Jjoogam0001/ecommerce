define setup_env
    $(eval include .env)
    $(eval export)
endef

local-env:
	$(call setup_env)


start:
	docker-compose -f ./build/docker-compose.yml -p ecommerce-service up --build --detach

stop:
	docker-compose -f ./build/docker-compose.yml -p ecommerce-service down --remove-orphans

bounce:
	docker-compose -f ./build/docker-compose.yml -p ecommerce-service down --remove-orphans
	docker-compose -f ./build/docker-compose.yml -p ecommerce-service up --build --detach

start-local: local-env
	go run cmd/ecommerce/main.go

run-local: local-env
	./main

build-ecommerce:
	go build ./cmd/ecommerce/main.go 

test-sociable-api:
	go test ./tests/sociable/ecommerce -v

gen-swagger:
	swag init --parseDependency --parseInternal --parseDepth 1 -g api/api.go -o api/docs

# Common utilities
docker-clean:
	@docker system prune --volumes --force