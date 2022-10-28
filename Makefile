.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: initDB
initDB:
	docker run --name=evo-test -e POSTGRES_PASSWORD=postgres -p 5555:5432 -d postgres

.PHONY: swagInit
swagInit:
	swag init -g ./cmd/main.go

.PHONY: composeUpBuild
composeUpBuild:
	docker compose up --build server

