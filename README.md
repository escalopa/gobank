# gobank
A SAAS banking application built with Go, and PostgreSQL, To manage your bank account, and make transactions.

## [LIVE API DOCS](http://37.46.128.188/gobank/docs/index.html)


## How to run

First, you need to have `docker` and `docker-compose` installed on your machine.

First we need to set the env variables, so copy the default env files to `.env` and `.env.db` files.

```bash
cp .env.example .env
cp .env.db.example .env.db
```

Then, run the following command:

```bash
docker-compose up
```

## Features

The applicatoin uses `paseto` for authentication.

### User

- Create a user
- Login
- Renew access token
- Upadte user

### Account
- Create an account
- Get all accounts (Of the logged in user)
- Delete an account (Soft delete, Can be restored)
- Restore an account (After has been deleted)

### Transaction
- Create a transaction
- Get all transactions (Of a specific account)

## Tech Stack

- Gin
- PostgresSQL
- gomock

## GRPC Services

The project uses GRPC besides the REST API, to communicate with db. But the GRPC are not implemented yet fully as the API.
