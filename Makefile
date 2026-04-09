run:
	go run main.go

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

build:
	go build -o bin/auth-service ./main.go
