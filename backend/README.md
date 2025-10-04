# backend

Backend service to serve REST APIs submit voucher seat assignments for airline campaign.

## Features

- Submit voucher assignments endpoints.
- Create a new flights, seats, and vouchers.
- List of all flights, seats, and vouchers.

## Architecture Overview

## Pre-Requisites

- [Go](https://go.dev/doc/install)
- SQLite

## Getting Started

```shell
cp .env.example .env
go run .
```

## Configuration

```
PORT=
```

## Endpoints

Create a new flights, to view just change the verb from `POST` to `GET`.

```shell
curl --location 'http://localhost:8080/api/v1/flights' \
--header 'Content-Type: application/json' \
--data '{
    "flight_numbers": ["GA33", "GA221"],
    "dept_date": "2025-10-04"
}'
```

Create a new seats, to view just change the verb from `POST` to `GET`.

```shell
curl --location 'http://localhost:8080/api/v1/seats' \
--header 'Content-Type: application/json' \
--data '{
 "flight_id": 23,
 "cabin": "BUSINESS",
 "labels": ["1A", "1B", "1C"]   
}'
```

Create a new vouchers, to view just change the verb from `POST` to `GET`.

```shell
curl --location 'http://localhost:8080/api/v1/vouchers' \
--header 'Content-Type: application/json' \
--data '{
    "code": "V2025X2",
    "flight_id": 23,
    "cabin": "ECONOMY"
}'
```
