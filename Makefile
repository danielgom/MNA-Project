.PHONY: run start-db
# Run the api
run:
	@docker build -t mna-project .
	@docker run -p 8080:8080 mna-project


start-db:
	@docker start MNA-postgres 2>/dev/null || docker run --name MNA-postgres -p 5432:5432 -e POSTGRES_PASSWORD=mypass123 -d postgres:16.0
	@docker cp schema/Initial_tables.sql MNA-postgres:/Initial_tables.sql
	@sleep 2
	@docker exec -t MNA-postgres psql -U postgres -f Initial_tables.sql

swagger-doc:
	swag init -g ./cmd/api/main.go --pd