# SQL Data access layer Setup instructions

## Pre-requisites
1. Docker
2. Postgres image installed on Docker

## Initial Database Setup instructions
Create a `.env` file inside the `sql_dal` folder and add the following contents
```
POSTGRES_DB=<your-database-name>
POSTGRES_USER=<your-database-user>
POSTGRES_PASSWORD=<your-database-password>
```

Open a terminal session inside `sql_dal` folder and Run
```
docker compose up -d
```

Check if the volume and the container is created then you are gtg.

Connect to the database via any db connection application

## How to connect sql_dal App to Database
In the `.env` file all the following contents
```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=<your-database-name>
DB_USER=<your-database-user>
DB_PASSWORD=<your-database-password>
DB_SSLMODE=disable
```

`DB_HOST=localhost` is for local testing only.

## Code Strucutre
* The `config/` folder contains the code for establishing the database connection, and migrating datamodels into the database.
* Add your datamodels inside the `data_models/` folders and keep the database queries inside the data model file itself.
* Once you have added the data models then make changes in the `config/migrations.go` file to ensure the data models are migrated
* Run the application and ensure that the database is migrated with the latest changes properly.
