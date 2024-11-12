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

Setup the table

```sql
CREATE TABLE cars(
  id SERIAL PRIMARY KEY, 
  brand VARCHAR(255) NOT NULL,
  model VARCHAR(255) NOT NULL,
  year INT NOT NULL,
  state VARCHAR(255) NOT NULL,
  color VARCHAR(255) NOT NULL,
  fuel_type VARCHAR(50) NOT NULL,
  body_type VARCHAR(50) NOT NULL
);
```

Seed the database

```sql
INSERT INTO cars(brand, model, YEAR, state, color, fuel_type, body_type)
SELECT (CASE FLOOR(RANDOM() * 5)::INT
            WHEN 0 THEN 'Toyota'
            WHEN 1 THEN 'Ford'
            WHEN 2 THEN 'Honda'
            WHEN 3 THEN 'BMW'
            WHEN 4 THEN 'Tesla'
        END) AS brand,
       (CASE FLOOR(RANDOM() * 5)::INT
            WHEN 0 THEN 'Camry'
            WHEN 1 THEN 'F-150'
            WHEN 2 THEN 'Civic'
            WHEN 3 THEN '3 Series'
            WHEN 4 THEN 'Model S'
        END) AS model,
       (CASE FLOOR(RANDOM() * 5)::INT
            WHEN 0 THEN 2024
            WHEN 1 THEN 2023
            WHEN 2 THEN 2022
            WHEN 3 THEN 2021
            WHEN 4 THEN 2020
        END) AS YEAR,
       (CASE FLOOR(RANDOM() * 3)::INT
            WHEN 0 THEN 'Operational'
            WHEN 1 THEN 'Under maintenance'
            WHEN 2 THEN 'Totalled'
        END) AS state,
       (CASE FLOOR(RANDOM() * 5)::INT
            WHEN 0 THEN 'Red'
            WHEN 1 THEN 'Green'
            WHEN 2 THEN 'Blue'
            WHEN 3 THEN 'Black'
            WHEN 4 THEN 'White'
        END) AS color,
       (CASE FLOOR(RANDOM() * 3)::INT
            WHEN 0 THEN 'Diesel'
            WHEN 1 THEN 'Petrol'
            WHEN 2 THEN 'Electric'
        END) AS fuel_type,
       (CASE FLOOR(RANDOM() * 3)::INT
            WHEN 0 THEN 'Sedan'
            WHEN 1 THEN 'SUV'
            WHEN 2 THEN 'Hatchback'
        END) AS body_type
FROM GENERATE_SERIES(1, 10000000) seq;
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
