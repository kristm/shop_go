DB := shop.db
TEST_DB := test.db
migrate_up: 
	migrate -path=internal/database/migrations -database "sqlite3://internal/database/$(DB)" -verbose up $(N)

migrate_up_test: 
	migrate -path=internal/database/migrations -database "sqlite3://internal/database/$(TEST_DB)" -verbose up $(N)

migrate_down: 
	migrate -path=internal/database/migrations -database "sqlite3://internal/database/$(DB)" -verbose down $(N)

migrate_down_test: 
	migrate -path=internal/database/migrations -database "sqlite3://internal/database/$(TEST_DB)" -verbose down $(N)

V ?= 1
migrate_reset: 
ifndef V
	migrate -path=internal/database/migrations -database "sqlite3://internal/database/$(DB)" -verbose force $(V)
else
	migrate -path=internal/database/migrations -database "sqlite3://internal/database/$(DB)" -verbose force $(V)
endif

migrate_version: 
	migrate -version -path=internal/database/migrations -database=sqlite3://internal/database/$(DB)

migrate_drop:
	migrate drop -f -ext=sql -database "sqlite3://internal/database/$(DB)"

migrate_create:
	migrate create -ext=sql -seq -dir=internal/database/migrations $(name) 

# does not work
migrate_goto:
	migrate goto $(version) -path=internal/database/migrations -database "sqlite3://internal/database/$(DB)"

clean_test:
	go clean -testcache

build_docker:
	docker build --progress=plain -t shop_go .

shell:
	docker run --entrypoint /bin/sh -it shop_go

run_shop:
	docker run --rm -p 8080:8080 shop_go

build:
	GIN_MODE=release go build cmd/main.go

vulcheck:
	govulncheck ./...
