# Testing SQL Builder

In this repo, I test my idea on creating a dynamic REST API.
I explore dynamic filtering of the output and the query based on query params.

## Setup

1. Database - Postgres

Setup the database and user for this project using the following SQL

```sql
CREATE DATABASE sqlbuildertest;
CREATE USER sqlbuilderuser WITH ENCRYPTED PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE sqlbuildertest TO sqlbuilderuser;
ALTER DATABASE sqlbuildertest OWNER TO sqlbuilderuser;
```

2. Application - Golang

- Clone the project

- Install packages

```bash
go mod tidy
```

- Run the application

```bash
go run main.go
```
