generate-protobufs:
	buf mod update
	buf generate

generate-sql:
	docker run --rm -v $PWD:/src -w /src sqlc/sqlc generate

run-migrations:
	go run cmd/migrate/main.go -database pricing -path migrations/pricing
	go run cmd/migrate/main.go -database catalog -path migrations/catalog
	go run cmd/migrate/main.go -database bookings -path migrations/bookings
	go run cmd/migrate/main.go -database tee_times -path migrations/tee_time

generate-data:
	go run cmd/generate/main.go
