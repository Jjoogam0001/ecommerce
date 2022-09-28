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
	go run cmd/generatebets/main.go

run-local: local-env
	./main

build-ecommerce:
	go build ./cmd/ecommerce/main.go -ldflags="-X main.BuildVersion=$(BUILD_VERSION_ARG)" ./cmd/ecommerce/main.go

test-sociable-api:
	go test ./test/sociable/ecommerce -v

gen-mocks-common:
	GO111MODULE=on mockgen --build_flags=--mod=mod -destination=./test/sociable/common/mocks/offer/offerclient.go -package mocks dev.azure.com/Derivco/Sports-CoreBetting/_git/bet-offer-client.git/offer/offerclient/offeroperations ClientService
	GO111MODULE=on mockgen --build_flags=--mod=mod -destination=./test/sociable/common/mocks/db/ecommercedb.go -package mocks -source=./internal/app/ecommerce/dbaccess/ecommercedb.go

gen-swagger:
	swag init -g internal/app/ecommerce/httpserver.go -o api/ecommerce/docs

# Common utilities
docker-clean:
	@docker system prune --volumes --force