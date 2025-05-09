build:
	docker compose up -d --build

start:
	docker compose up -d

stop:
	docker compose down

logs/app:
	docker compose logs -f --no-log-prefix app

migrate/create:
	docker run -v ./migrations:/migrations --network host migrate/migrate create -ext sql -dir ./migrations $(name)

swagger:
	docker run --rm -v ./:/code ghcr.io/swaggo/swag:latest init

swagger-mac:
	docker run --platform linux/amd64 --rm -v ./:/code ghcr.io/swaggo/swag:latest init
