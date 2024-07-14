
.DEFAULT_GOAL := run


MAIN_GO := main.go


run:
	CompileDaemon --build="go build $(MAIN_GO)" --command="./main" 


migrate:
	go run api/migrate/migrate.go

up:
	docker-compose up -d

down:
	docker-compose down

