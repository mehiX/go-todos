# TODO list (coding challenge)

Write a simple REST API for a Todo list application. 

- Develop and test it without opening the browser or using curl, postman, etc.
- Use as few external dependencies as possible

## Requirements

- go >= 1.19
- Docker - to run it as container

## General considerations

I chose a flat structure since the application is small and there are not too many pieces. If it keeps growing I would consider using different packages for service, repositories, domain, etc.

I use [go-chi](https://github.com/go-chi/chi) because it is elegant and it saves me time to parse path parameters and test on http methods. 

The application uses 2 types of persistence: SQL database (MariaDB/Mysql) or in-memory (if connecting to the database fails).

At startup the application tries to connect to the database. It includes a retry mechanism with exponential backoff. The retry parameters are hardcoded for now, should be passed in as arguments.

If an empty string is passed as `dsn` argument then the application will not even try to connect to a database and use the In Memory persistence directly.

There are integration tests provided in `/tests`.

There are some unit tests provided, but the coverage is not great due to time limitations:

```shell
go test -v -cover -coverprofile=cover.out ./pkg/todos/...
go tool cover -html=cover.out
```

## Endpoints provided

The list of provided endpoints is dynamically generated (no update needed if a new pattern + handler is added) and printed when the program starts.

```
[GET]           /health


[GET]           /todos/
[POST]          /todos/

[GET]           /todos/search/tags


[DELETE]        /todos/{id:[0-9a-zA-Z-]+}/
[GET]           /todos/{id:[0-9a-zA-Z-]+}/
[PUT]           /todos/{id:[0-9a-zA-Z-]+}/
```


## Build and run

```shell
go build -o todos ./main.go
./todos --help
```

Run the application in a shell:

```shell
./todos --http 127.0.0.1:7070 --dsn ''
```

and the integration tests in a separate shell:

```shell
go test ./tests/... -args -api-url='http://127.0.0.1:7070'
```

## Run with Docker

Build and run the container and keep attached to see the logs

```shell
docker compose up --build
```

Run the tests against the running container

```shell
go test ./tests/... -args -api-url='http://127.0.0.1:7070'
```

