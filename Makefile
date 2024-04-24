migrate:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5432/orderswb?sslmode=disable' up

migrate_down:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5432/orderswb?sslmode=disable' down