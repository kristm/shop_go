migrate_up: 
	migrate -path=internal/database/migrations -database "sqlite3://internal/database/shop.db" -verbose up $(N)

migrate_down: 
	migrate -path=internal/database/migrations -database "sqlite3://internal/database/shop.db" -verbose down $(N)

migrate_reset: 
	migrate -path=internal/database/migrations -database "sqlite3://internal/database/shop.db" -verbose force $(V)

migrate_version: 
	migrate -version -path=internal/database/migrations -database=sqlite3://internal/database/shop.db

migrate_drop:
	migrate drop -f -ext=sql -database "sqlite3://internal/database/shop.db"

migrate_create:
	migrate create -ext=sql -seq -dir=internal/database/migrations $(name) 

# does not work
migrate_goto:
	migrate goto $(version) -path=internal/database/migrations -database "sqlite3://internal/database/shop.db"

clean_test:
	go clean -testcache

build_docker:
	docker build --progress=plain -t shop_go .

shell:
	docker run --entrypoint /bin/sh -it shop_go

run_shop:
	docker run --rm -p 8080:8080 shop_go
