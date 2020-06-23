DB_URL=${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

build:
	@go build -o bin/auth -v
run:
	@export $(cat .env | xargs) && ./bin/auth
migrate:
	@migrate --path=db/migrations --database=${DB_URL} up
drop:
	@migrate --path=db/migrations --database=${DB_URL} down
clean:
	@go clean -o bin/auth
vendor:
	@go mod vendor