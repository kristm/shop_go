migrate_up: 
	migrate -path=internal/database/migrations -database "sqlite3://shop_test.DB" -verbose up $(N)

migrate_down: 
	migrate -path=internal/database/migrations -database "sqlite3://shop_test.DB" -verbose down $(N)

migrate_reset: 
	migrate -path=internal/database/migrations -database "sqlite3://shop_test.DB" -verbose force $(V)

migrate_version: 
	migrate -version -path=internal/database/migrations -database=sqlite3://shop_test.DB

migrate_drop:
	migrate drop -f -ext=sql -database "sqlite3://shop_test.DB"

migrate_create:
	migrate create -ext=sql -seq -dir=internal/database/migrations $(name) 

# does not work
migrate_goto:
	migrate goto $(version) -path=internal/database/migrations -database "sqlite3://shop_test.DB"

clean_test:
	go clean -testcache
