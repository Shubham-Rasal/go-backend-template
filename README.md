ackend System Development

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


Don't forget to add the migration commands to the Makefile:

```Makefile

migrateup:
	./migrate -path db/migration -database "postgresql://root:password@localhost:5432/blog?sslmode=disable" -verbose up ($version)

migratedown:
	./migrate -path db/migration -database "postgresql://root:password@localhost:5432/blog?sslmode=disable" -verbose down ($version)

.PHONY: createdb dropdb creatpg runpg  stoppg migrateup migratedown
```

TODO: Enable ssl on postgress

## Generate code from SQL

Following is a explaination for choosing sqlc for interacting with db.

Adapted from the blog - https://blog.jetbrains.com/go/2023/04/27/comparing-db-packages/

Blog Title: "Comparing database/sql, GORM, sqlx, and sqlc" by Sergey Kozlovskiy

- Comparison of Go packages for working with databases: `database/sql`, `GORM`, `sqlx`, and `sqlc`.
- `database/sql`: Standard library package for database operations in Go.
- `sqlx`: Extension of `database/sql` with features like named parameters and struct scanning.
- `sqlc`: SQL compiler generating type-safe code for raw SQL queries.
- `GORM`: Full-featured Go ORM library for advanced querying.
- Comparison factors: Features, ease of use, performance, and speed.
- Code examples provided for each package.
- Performance benchmarks show GORM excels for small record counts but lags for large records.
- Conclusion: Choose the package based on your specific needs as a developer.
- GORM for advanced querying and clean code.
- `database/sql` and `sqlx` for basic queries.
- `sqlc` for backend developers with many queries and tight deadlines.

Choose `sqlc` for my Go database needs because it generates type-safe code, boosts productivity, offers good performance, is easy to use, has community support, and is ideal for backend development. It simplifies database interactions, making it a valuable choice for Go developers.


### sqlc

sqlc is a Go library that generates type-safe code from SQL queries. 

Installation instructions can be found [here](https://docs.sqlc.dev/en/latest/overview/install.html).

First create a sqlc.yaml file in the root directory of the project. This file will contain the configuration for sqlc. Following is the configuration for our blog application.

```yaml

version: "1"
packages:
  - name: "db"
    path: "./db/sqlc"
    queries: "./db/query/"
    schema: "./db/migration/"
    engine: "postgresql"
    emit_json_tags: true
    emit_prepared_queries: true
    emit_interface: false
    emit_exact_table_names: false

```

- `name`: Name of the package.
- `path`: Path to the package.
- `queries`: Path to the directory containing the queries.
- `schema`: Path to the directory containing the schema.
- `engine`: Database engine.
- `emit_json_tags`: Emit json tags for the generated code.
- `emit_prepared_queries`: Emit prepared queries for the generated code.
- `emit_interface`: Emit interface for the generated code.
- `emit_exact_table_names`: Emit exact table names for the generated code.

Now we can run the following command to generate the code.

```bash
sqlc generate
```

This will give error because we have not created the queries yet. We can create the queries in the `db/query` directory. We can create a file called `user.sql` in the `db/query` directory and add the following query.

```sql

-- name: CreateUser :one
INSERT INTO users (
	username,
	password_hash,
	email
) VALUES (
	$1,
	$2,
	$3
) RETURNING *;

...

```

This query will create a user in the database. The `:one` in the query indicates that this query will return one row. The `$1`, `$2`, `$3` are the parameters for the query. The parameters are used to prevent SQL injection attacks.

Now we can run the following command to generate the code.

```bash
sqlc generate
```

This will generate the code in the `db/sqlc` directory. We can use this code to interact with the database.

Three files generated:

- `db.go`: Contains the database connection code.
- `models.go`: Contains the models for the tables.
- `users.sql.go`: Contains the code for the queries.

## Testing the generated code

Tests in Go are written using the testing package. You can read more about it [here](https://golang.org/pkg/testing/).

We can create a file called `main_test.go` in the `db/sqlc` directory and add the required test.
User specific tests are added in `user_test.go` file.

For example, to test the CreateUser query, we can add the following test.

```go

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		Username: util.RandomUserName(),
		Role:     util.RandomRole(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Role, user.Role)
}

```

Note: The utils file contains a file `random.go` to generate random strings and numbers.


## Continuous Integration

Continuous Integration is the practice of merging all developer working copies to a shared mainline several times a day. It was first introduced by Grady Booch in his 1991 method, although he did not advocate integrating several times a day. 

We will use Github Actions for CI. You can read more about Github Actions [here](https://docs.github.com/en/actions).

One repetitive task which I can identify for now is running the tests. We can add a Github Action to run the tests on every push to the repository. This will automate the process of running the tests.

We can create a file called `ci.yml` in the `.github/workflows` directory and add the following code.

```yaml
# FILEPATH: /d:/coding/goPrograms/practice/.github/workflows/ci.yml
# This workflow will build a golang project and run tests.
# It uses PostgreSQL as a service and runs migrations before running tests.
# For more information see: https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go

name: ci-test

on:
  push:
    branches: [ "master" ]
  pull_request:
    branches: [ "master" ]

jobs:

  test:
    name: Test
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-bookworm
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_USER: root
          POSTGRES_DB: blog
        options: >-
          --health-cmd pg_isready
          --health-interval 5s
          --health-timeout 3s
          --health-retries 3
        ports:
          # Maps tcp port 5432 on service container to the host
          - 5432:5432
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '^1.20'
    
    - name: Install Migration Tool
      run: |
        curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
        mv migrate /usr/local/bin/
        which migrate
    
    - name: Run migrations
      run: make migrateup step=3

    - name: Test
      run: make test

```
This specific workflow is triggered on two events: when a push is made to the "master" branch or when a pull request is made to the "master" branch.

The workflow consists of a single job named "test" that runs on the latest version of Ubuntu. This job has several steps and also uses a service, which is a Docker container that the job can interact with.

The service used here is a PostgreSQL database. The image used for this service is postgres:16-bookworm. The environment variables POSTGRES_PASSWORD, POSTGRES_USER, and POSTGRES_DB are set for this service. The service also has health check options and a port mapping.

The steps in the job include:

Checking out the repository using the actions/checkout@v3 action.
Setting up Go using the actions/setup-go@v4 action with a version specification of '^1.20'.
Installing a migration tool by downloading it from a URL, moving it to /usr/local/bin/, and verifying its installation with which migrate.
Running migrations using a make command.
Running tests using another make command.




## Rest API using go fiber

To create a rest api, we will use a library called [fiber](https://docs.gofiber.io/). It is a powerful wrapper around the go standard net/http package with many built-in features.

You can easily create routes and handlers using fiber as follows.

```go

package main

import "github.com/gofiber/fiber/v2"

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Listen(":3000")
}

```

All server related stuff is going in the api diretory. A `server` struct encapsulates all the features of our backend using structure embedding.


```go

type Server struct {
	store     db.Store
	router    *fiber.App
	validator *validator.Validate
}

//Install it using : go get github.com/go-playground/validator/v10

```

We use a library called [validator](https://github.com/go-playground/validator) for input validation.
It uses tags to figure out the correct field in input and output params.







