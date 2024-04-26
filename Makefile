migrate:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5432/orderswb?sslmode=disable' up

migrate_down:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5432/orderswb?sslmode=disable' down

run_server:
	go run cmd/main.go

run_script:
	go run cmd/nats-streaming-script/main.go

nats:
	nats-streaming-server