include .env

create_migration:
	migrate create -ext=sql -dir=internal/database/migrations -seq $(name)

migrate_up:
	migrate -path=internal/database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_ADDRESS}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=internal/database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_ADDRESS}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

migrate_force:
	migrate -path=internal/database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_ADDRESS}:${DB_PORT}/${DB_NAME}?sslmode=disable" force $(version)

migrate_rollback:
	migrate -path=internal/database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_ADDRESS}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down $(version)