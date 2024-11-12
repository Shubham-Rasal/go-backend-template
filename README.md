# Backend System Development

## Table of Contents

- [Backend System Development](#backend-system-development)
	- [Table of Contents](#table-of-contents)
	- [Database Schema Generation (Postrgesql)](#database-schema-generation-postrgesql)
	- [Getting a Postgresql Database (Docker)](#getting-a-postgresql-database-docker)
	- [Makefile](#makefile)
	- [TablePlus](#tableplus)
	- [Database Schema Migration](#database-schema-migration)
	- [Generate code from SQL](#generate-code-from-sql)
		- [sqlc](#sqlc)
	- [Testing the generated code](#testing-the-generated-code)
	- [Continuous Integration](#continuous-integration)
	- [Rest API using go fiber](#rest-api-using-go-fiber)
		- [Routes and Handlers](#routes-and-handlers)
		- [Authentication using JWT and PASETO](#authentication-using-jwt-and-paseto)
			- [JWT](#jwt)
			- [PASETO](#paseto)
		- [Validation](#validation)
	- [Environment variables](#environment-variables)
		- [Middleware](#middleware)
	- [Mock DB](#mock-db)
	- [Dockerisation](#dockerisation)
		- [Why Docker?](#why-docker)
		- [Multistage Build](#multistage-build)
		- [Build the image](#build-the-image)
		- [Docker networking](#docker-networking)
		- [Further build optimisations](#further-build-optimisations)
	- [Using Docker Compose](#using-docker-compose)

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
docker run --name some-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14-bookworm
```

What does this command do?

- `docker run`: This command is used to run a docker container.
- `--name some-postgres`: This is the name of the container. You can use any name you want.
- `-p 5432:5432`: This is used to map the port 5432 of the container to the port 5432 of the host machine. Postgresql runs on port 5432 by default.
- `-e` is used to set environment variables. We are setting the username and password for the database.
- `-d` is used to run the container in the background.
- `postgres:14-bookworm` is the name of the image we want to run. This is the image for Postgresql version 16-bookworm.

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
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14-bookworm

```

Note: You also get auto-completion for Makefiles in your terminal.

For now we can add the following commands to the Makefile.

```Makefile

creatpg:
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14-bookworm

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
migrate create -ext sql -dir db/migration -seq init_schema -verbose
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
        image: postgres:14-bookworm
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

The service used here is a PostgreSQL database. The image used for this service is postgres:14-bookworm. The environment variables POSTGRES_PASSWORD, POSTGRES_USER, and POSTGRES_DB are set for this service. The service also has health check options and a port mapping.

The steps in the job include:

Checking out the repository using the actions/checkout@v3 action.
Setting up Go using the actions/setup-go@v4 action with a version specification of '^1.20'.
Installing a migration tool by downloading it from a URL, moving it to /usr/local/bin/, and verifying its installation with which migrate.
Running migrations using a make command.
Running tests using another make command.




## Rest API using go fiber

To create a rest api, we will use a library called [fiber](https://docs.gofiber.io/). It is a powerful wrapper around the go standard net/http package with many built-in features.


### Routes and Handlers

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

The `store` field is used to interact with the database. The `router` field is used to create routes and handlers. The `validator` field is used to validate the input and output parameters.



NOTE: To test the APIs, following tools can be used

- postman
- thunderclient (vscode extension)
- cURL

### Authentication using JWT and PASETO

You can read more about JWT [here](https://jwt.io/) and PASETO [here](https://paseto.io/).

#### JWT

JWT stands for JSON Web Token. It is an open standard that defines a compact and self-contained way for securely transmitting information between parties as a JSON object. This information can be verified and trusted because it is digitally signed.

JWTs can be signed using a secret (with the HMAC algorithm) or a public/private key pair using RSA or ECDSA. 

#### PASETO

PASETO stands for Platform-Agnostic Security Tokens. It is a specification and reference implementation for secure stateless tokens. It is a more secure alternative to JWT.

PASETO is a secure alternative to the JOSE standards (JWT, JWE, JWS) and can be used in a wide variety of use cases. It is a secure alternative to JWT.

I have implementated both using a common interface `Maker` which has the following features:

```go

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

```

Any struct which implements this interface can be used to create and verify tokens. The `CreateToken` method is used to create a token for a given username and duration. The `VerifyToken` method is used to verify a token and return the payload.

The `Maker` interface is implemented by two structs, `JWTMaker` and `PasetoMaker`. The `JWTMaker` struct is used to create and verify JWT tokens. The `PasetoMaker` struct is used to create and verify PASETO tokens.

```go

type JWTMaker struct {
	secretKey string
}

type PasetoMaker struct {
	pasetoKey string
}

```

The `secretKey` field is used to sign the JWT token. The `pasetoKey` field is used to sign the PASETO token.

Due to simplicity of paseto and its security features, I have used it for authentication.

NOTE: It can be easily replaced by JWT by changing the `Maker` interface and the `NewPasetoMaker` function.


### Validation

We use a library called [validator](https://github.com/go-playground/validator) for input validation.
It uses tags to figure out the correct field in input and output params.

```go

type ListUsersQueryParams struct {
	PageId int32 `validate:"required,min=1"`
	PageSize int32 `validate:"required,min=5,max=10"`
}

func (server *Server) listUsers(c *fiber.Ctx) error {
	var params ListUsersQueryParams
	err := c.QueryParser(&params)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	// Validate request
	if err := server.validator.Struct(params); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}
  ...
```


Here we are validating the query params for the list users endpoint. The `validate` tag is used to specify the validation rules. The `required` tag is used to specify that the field is required. The `min` and `max` tags are used to specify the minimum and maximum values for the field.

## Environment variables

All global variables like `serverAddress` , `dbSource` are replace by environment variables.
These are loaded using a utility called [viper](https://github.com/spf13/viper).

One cool feature I like about the viper setup is that the env variables are loaded into a struct which makes it really easy for auto-complete and interllisense.

### Middleware

Middlewares are used to intercept requests and responses. You can read more about it [here](https://docs.gofiber.io/guide/middleware).

You can add various features using middleware like logging, authentication, tracing, etc.

Since we are using fiber, the middleware creation is super easy. We can create a middleware for authentication as follows.

```go

func Protected(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		//do your auth or any other stuff 
		return ctx.Next()
	}
}
```
 The auth middleware helps you to set the user in the context of the request. This can be used to get the user in the handler. This is done by using the `ctx.Locals` method.

 ```go

	// some stuff
	user := ctx.Locals("user").(*jwt.Token)
// other stuff

 ```

 This can used in the handlers to check if this user is allowed to perform the action or not.
This is called authorization.


## Mock DB

There are two ways to mock the database.

- using a mock in memory database, but requires a lot of code
- using db stubs (mocks) using [mock](https://github.com/uber-go/mock).


For mock to work in reflect mode (dev mode) we need to add a black import to the `main.go` file.

For making a mock store we need a generalised interfact `Store`. We abstract it away from out sql postgress store by removing the methods to an interface and make the old store as `SQLStore`.

```go

//generic interface for all the queries
//anyone who wants to be a store must implement this interface
type Store interface {
	Querier //genrated by sqlc
	LikeTx(ctx context.Context, arg LikePostParams) error
}

// has db and set of queries to interact with the database
type SQLStore struct {
	db *sql.DB
	*Queries
}

// constructor
func NewStore(db *sql.DB) Store {

	//sqlstore implements the store interface as the queries field is a part of it
	//and the likeTx method is implemented by the sqlstore
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

```

We can add all the funcion signatures to the store manually but a better way would be to somehow generate the signatures automatically. This can be easily done by toggling the following parameter in`sqlc.yaml`

```yml

# version: "1"
# packages:
#   - name: "db"
#     path: "./db/sqlc"
#     queries: "./db/query/"
#     schema: "./db/migration/"
#     engine: "postgresql"
#     emit_json_tags: true
#     emit_prepared_queries: true
    emit_interface: true
#     emit_exact_table_names: false
#     emit_empty_slices: true

```
This generates a `Querier` strunt to embed in the store interface.

The `SQLStore` already implements the interfact `Store`. Now we use this definition to generate a mock store using the following command -

```bash
mockgen -package mockdb -destination db/mock/store.go <insert-your-package-name>/db/sqlc Store
```

This generates a store.go in `mock` folder which can now be used to test apis without accessing the actual db.

For example,

```go

{
			name:   "OK",
			userId: int64(user.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatch(t, user, response)
			},
		},
....

```

In software testing, a mock object is a simulated object that mimics the behavior of the real object in a controlled way. A mock object returns a specific pre-programmed response to a method call.

A stub is a type of mock object that returns a pre-determined response to a method call. Stubs are used to replace real objects in a test and provide a predictable response to the code being tested.

In the provided code, `buildStubs` is a function that creates a mock object of the `store` and sets an expectation that the `GetUser` method will be called once with the `user.ID` as an argument. This mock object is used to test the getUserHandler function in isolation from the actual database.In software testing, a mock object is a simulated object that mimics the behavior of the real object in a controlled way. A mock object returns a specific pre-programmed response to a method call.



## Dockerisation

![Docker Multistage Build](image-1.png) 

### Why Docker?

Docker is a tool designed to make it easier to create, deploy, and run applications by using containers. Containers allow a developer to package up an application with all of the parts it needs, such as libraries and other dependencies, and deploy it as one package.

### Multistage Build

Since golang compiler and all the dependencies are required to run the application, we need to create a docker image with all the dependencies. But this image will be very large in size. So we use a multistage build to create a smaller image.

The first stage is used to build the application. The second stage is used to run the application. The second stage does not contain the compiler and the dependencies. It only contains the binary and the required dependencies.

The Dockerfile for the application is as follows.

```Dockerfile

#Build stage
FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

#Final stage
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .

# cpy the app.env file to the container
COPY --from=builder /app/app.env .

EXPOSE 8080

CMD ["/app/main"]

```

The first stage uses the golang:1.21.4-alpine3.18 image. This image contains the golang compiler and the dependencies required to build the application. The WORKDIR instruction sets the working directory to /app. The COPY instruction copies all the files from the current directory to the /app directory in the container. The RUN instruction runs the go build command to build the application. The -o flag is used to specify the name of the output file. The CMD instruction specifies the command to run when the container starts.

### Build the image

To build the image, we can run the following command.

```bash

docker build -t blog-backend:latest .

```

The -t flag is used to specify the name of the image. The . specifies the current directory.

Note: If the image already exists, you can use the --no-cache flag to build the image from scratch.

```bash

docker build --no-cache -t blog-backend:latest .

```



### Docker networking

Notes on docker networking can be found [here](Notes.md/#docker-network).

First create a network `blognet` -

```bash

docker network create blognet

```

Modify the creatpg command in the Makefile to add the network flag.

```Makefile

createpg:
	docker run --name some-postgres --network blognet -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14-bookworm

```

While making a container out of the image, we need to add the network flag.

```bash

docker run --network blognet --name blog-backend -p 8080:8080 -e DB_SOURCE="postgresql://root:password@some-postgres:5432/blog?sslmode=disable" blog-backend:latest

```

Note: since both the containers are on the same network, we can use the container name as the host name in the connection string.


### Further build optimisations

Found a good method for speeding up builds in golang projects. It uses a cache layer to store the dependencies and the build artifacts. This reduces the build time significantly. 
You can find the article [here](https://www.docker.com/blog/containerize-your-go-developer-environment-part-2/).

Another one is [here](https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324).

```Dockerfile

# Replacing this in the build stage
# RUN go build -o main main.go

# with this

WORKDIR /app

COPY . .

ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main main.go

```

This will cache the dependencies and the build artifacts. This will reduce the build time significantly.


## Using Docker Compose

All the services 