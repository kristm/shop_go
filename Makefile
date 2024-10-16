migrate_up: 
	migrate -path=internal/database/migrations -database "sqlite3://shop_test.DB" -verbose up

migrate_down: 
	migrate -path=internal/database/migrations -database "sqlite3://shop_test.DB" -verbose down

migrate_reset: 
	migrate -path=internal/database/migrations -database "sqlite3://shop_test.DB" -verbose force 1

migrate_drop:
	migrate drop -f -ext=sql -database "sqlite3://shop_test.DB"

migrate_create:
	migrate create -ext=sql -seq -dir=internal/database/migrations $(name) 

# does not work
migrate_goto:
	migrate goto $(version) -path=internal/database/migrations -database "sqlite3://shop_test.DB"
