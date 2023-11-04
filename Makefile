creatpg:
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-bookworm

runpg:
	docker start some-postgres && docker ps

stoppg:
	docker stop some-postgres

createdb:
	docker exec -it some-postgres createdb --username=root --owner=root blog

dropdb:
	docker exec -it some-postgres dropdb blog 

migrateup:
	./migrate -path db/migration -database "postgresql://root:password@localhost:5432/blog?sslmode=disable" -verbose up

migratedown:
	./migrate -path db/migration -database "postgresql://root:password@localhost:5432/blog?sslmode=disable" -verbose down

.PHONY: createdb dropdb creatpg runpg  stoppg migrateup migratedown