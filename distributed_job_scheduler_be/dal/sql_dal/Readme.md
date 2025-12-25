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