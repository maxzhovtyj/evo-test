.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: initDB
initDB:
	docker run --name=evo-test -e POSTGRES_PASSWORD=postgres -p 5555:5432 -d postgres