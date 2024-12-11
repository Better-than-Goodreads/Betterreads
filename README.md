# Betterreads API

The Betterreads API serves as a backend prototype inspired by Goodreads. Its goal is to provide a platform where users can share opinions about books and rate them.

## Dependencies

The project is designed to run in Docker with a local database. Therefore, the only prerequisite is having Docker installed on your system.

### Environment Variables

To set up the environment, ensure the `.env` file exists in the `src` directory with the following variables:

```shell
ENVIRONMENT=development
PORT=port
HOST=host
DATABASE_HOST=host
DATABASE_PORT=port
DATABASE_NAME=name
DATABASE_USER=user
DATABASE_PASSWORD=password
JWT_SECRET=any
JWT_DURATION_HOURS=1
```

Additionally, another `.env` file is required inside the `/database` directory:

```shell
POSTGRES_USER=user
POSTGRES_PASSWORD=user123
POSTGRES_DB=db
```

> **Note:** The values for environment variables may vary depending on the database configuration. It's crucial that `user`, `password`, and `database_name`/`postgres_db` match in both files.

### Starting the Project and Database

Start the project and database with Docker Compose:

```shell
docker compose up
```

## Documentation

The documentation is automated using Swagger and Swag for Go. To generate the documentation, install the Swag CLI with:

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

Then, run the following command inside the `src/` directory to generate the documentation:

```shell
swag init -g cmd/main.go
```

To view the documentation, with the server running, navigate to:
[Swagger UI](http://localhost:8080/swagger/index.html#/)

Alternatively, replace `8080` with your configured port:

```
http://localhost:PORT/swagger/index.html#/
