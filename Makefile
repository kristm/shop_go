migrate_up: 
	migrate -path=internal/database/migrations -database "sqlite3://shop_test.DB" -verbose up

migrate_down: 
	migrate -path=internal/database/migrations -database "sqlite3://shop_test.DB" -verbose down

migrate_reset: 
	migrate -path=internal/database/migrations -database "sqlite3://shop_test.DB" -verbose force 1

migrate_create:
	migrate create -ext=sql -dir=internal/database/migrations $(name)
