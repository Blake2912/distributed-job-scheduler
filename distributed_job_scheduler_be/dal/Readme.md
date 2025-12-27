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

## How to connect sql_dal App to Database
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

## Code Strucutre/Standards to follow
* The `config/` folder contains the code for establishing the database connection, and migrating datamodels into the database.
* Add your datamodels inside the `sql_dal/data_models/` folders and keep the database queries inside the data model file itself.
* Once you have added the data models then make changes in the `sql_dal/config/migrations.go` file to ensure the data models are migrated
* Run the application and ensure that the database is migrated with the latest changes properly.
