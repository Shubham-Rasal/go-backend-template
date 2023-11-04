Backend System Development

## Database Schema Generation (Postrgesql)

Before starting development, we need to create a database schema. We can take help of tools like [dbdiagram.io](https://dbdiagram.io/home) to create a database schema.

![Database Schema](schema.png)

This is a simple database schema for a blog application. We have 2 tables, `users` and `posts`. A user can have multiple posts, but a post can only have one user. This is a one-to-many relationship.

Users can follow other users. This is a many-to-many relationship. We need a third table to store this relationship. We call this table `follows`. This table has two columns, `following_user_id` and `following_user_id`. Both of these columns are foreign keys to the `users` table.

Once created you can export the schema as a SQL file. This file can be used to create the database schema in Postgresql (or any other database).

## Getting a Postgresql Database (Docker)

We will use Docker to run Postgresql. You can install Docker from [here](https://docs.docker.com/get-docker/).

Once installed, you can run the following command to start a Postgresql database.

```bash
docker run --name some-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-bookworm
```

What does this command do?

- `docker run`: This command is used to run a docker container.
- `--name some-postgres`: This is the name of the container. You can use any name you want.
- `-p 5432:5432`: This is used to map the port 5432 of the container to the port 5432 of the host machine. Postgresql runs on port 5432 by default.
- `-e` is used to set environment variables. We are setting the username and password for the database.
- `-d` is used to run the container in the background.
- `postgres:16-bookworm` is the name of the image we want to run. This is the image for Postgresql version 16-bookworm.

After running this command, you can check if the container is running by running the following command.

```bash
docker ps
```

![Postgres Docker Container](postgress-docker-container.png)

You can connect to the database using the following command.

```bash
docker exec -it some-postgres psql -U root
```

This will open the Postgresql shell inside the container.

Note: You can use docker logs < container-id > to see the logs of the container.


## Makefile

To make running commands easier, we will use a Makefile. You can read more about Makefiles [here](https://www.gnu.org/software/make/manual/make.html). Makefiles are used to automate tasks.

For example to start the postgresql database, instead of running the longer command, we can add it to the Makefile and run it using `make creatpg`.

```Makefile
creatpg:
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-bookworm

```

Note: You also get auto-completion for Makefiles in your terminal.

For now we can add the following commands to the Makefile.

```Makefile

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


.PHONY: createdb dropdb creatpg runpg 

```

## TablePlus

You can use [TablePlus](https://tableplus.com/) to connect to the database. You can download it from [here](https://tableplus.com/).

Doing this will allow you to see the database schema in a GUI. You can also run queries from the GUI. It will make the development process easier.

![TablePlus](tableplus.png)

Use the exported sql file (schema.sql in my case) and copy all the commands and paste them in the SQL editor of tableplus and run all.

![Tables after running the command](tables.png)


## Database Schema Migration 

Migrations are used to manage the database schema. We can use migrations to create tables, add columns, remove columns, etc. We can also use migrations to seed the database with some data.

Migrations are essential for the development process. It allows us to make changes to the database schema without losing any data. For our blog applicaton, we will be using a library called [golang-migrate](https://github.com/golang-migrate/migrate) to manage our migrations.

Installation instructions can be found [here](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).


First we need to create a migration called init_schema. This migration will create the tables in the database. We can create this migration using the following command. Make a directory called `db/migration` and run the following command.

```bash
./migrate create -ext sql -dir db/migration -seq init_schema -verbose
```

This will create a file called `db/migration/000001_init_schema.up.sql` and `db/migration/000001_init_schema.down.sql`. The up file will contain the commands to create the tables and the down file will contain the commands to drop the tables. 

Initially the files are empty. We can copy the commands from the exported sql file and paste them in the up file. For down file, we can add the following commands.

```sql
DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;
```

Now to run the up migration, following command can be used:

Note: Make sure the postgres container is running.

```bash
./migrate -path db/migration -database "postgresql://root:password@localhost:5432/blog?sslmode=disable" -verbose up
```

Confirm that the tables are created by refreshing in tableplus.

Similarly rollback migration can be run like:

```bash
./migrate -path db/migration -database "postgresql://root:password@localhost:5432/blog?sslmode=disable" -verbose down
```

TODO: Enable ssl on postgress

Don't forget to add the migration commands to the Makefile:

```Makefile

migrateup:
	./migrate -path db/migration -database "postgresql://root:password@localhost:5432/blog?sslmode=disable" -verbose up

migratedown:
	./migrate -path db/migration -database "postgresql://root:password@localhost:5432/blog?sslmode=disable" -verbose down

.PHONY: createdb dropdb creatpg runpg  stoppg migrateup migratedown
```