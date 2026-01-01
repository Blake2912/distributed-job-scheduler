# Data access layer Setup instructions

## Pre-requisites
1. Docker

## Initial Database and Redis Setup instructions
Create a `.env` file inside the `dal` folder and add the following contents
```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=<your-database-name>
DB_USER=<your-database-user>
DB_PASSWORD=<your-database-password>
DB_SSLMODE=disable
```

Open a terminal session inside `dal` folder and Run for creating both postgres db and redis db
```
docker compose --profile redis --profile db up -d
```
The above command will pick the instructions from the `docker-compose.yml` file and create the database locally

Check if the volume and the container is created then you are gtg.

Connect to the database via any db connection application

## Running Database migrations
We need to keep the logic for running migrations separately because we will be scaling the dal service pods as per the demand and we don't want the migrations to be running on each pod spawn which will be costly and lock the database for sometime, so do the following:

In the `.env` file add the following contents if not added above
```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=<your-database-name>
DB_USER=<your-database-user>
DB_PASSWORD=<your-database-password>
DB_SSLMODE=disable
```

`DB_HOST=localhost` is for local testing only.
Run
```
go run ./cmd/migrate
```
Ensure the migrations are properly applied from the application logs and from the database end

## How to connect App to Database
In the `.env` file add the following contents if not added above
```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=<your-database-name>
DB_USER=<your-database-user>
DB_PASSWORD=<your-database-password>
DB_SSLMODE=disable
```

Run
```
go run ./cmd/server
```


## Code Strucutre/Standards to follow
* The `sql_dal/config/` folder contains the code for establishing the database connection, and migrating datamodels into the database.
* Add your datamodels inside the `sql_dal/models/` folder
* Once you have added the data models then make changes in the `sql_dal/config/migrations.go` file to ensure the data models are migrated
* Create separate repositories for each data model that was added inside the `sql_dal/repository` folder, here the main folder will contain the interface of the repo and the `postgres` folder should have the concrete implementation of the queries that we are trying to do to the database. NOTE: The database needs to be accessible only inside this layer, the service or the handler should not access the database directly.
* Create separate service inside `internal/service` to have the buisiness logic (if any)
* Create separate handlers (inside `internal/api/handlers`) and routes (inside `internal/api/routes`) for accessing the queries that was exposed as part of the service and repo.
* Ensure the created repo, service, handler are instantiated properly under the `internal/container/container.go` file (Required for dependency injection).
* Ensure the data flow is like `handler -> service -> repository -> model`
* Each handler should have its separate `<name>_routes.go` file, and register the routes inside the main `routes.go` file.
