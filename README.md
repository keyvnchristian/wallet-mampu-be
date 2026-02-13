# Wallet Service Mampu.IO (Golang + PostgreSQL)

A simple digital wallet service that supports:

-   Withdraw API
-   Balance Inquiry API
-   Health Check (DB Ping Additional)

Built with:

-   Go (net/http)
-   PostgreSQL
-   sqlx
-   Squirrel (SQL builder)

------------------------------------------------------------------------

## Architecture

    transport       → HTTP transport layer  
    usecase         → Business logic  
    repository      → Database access  
    config          → Configuration & DB connection  

Withdraw operation is concurrency-safe using:

-   Database transaction
-   `SELECT ... FOR UPDATE`
-   Row-level locking

------------------------------------------------------------------------

## Requirements

Install the following:

-   Go 1.20+\
    https://go.dev/dl/

-   PostgreSQL 13+\
    https://www.postgresql.org/download/

Optional: - curl - Postman

------------------------------------------------------------------------

## Setup

### 1. Clone Repository

``` bash
git clone <your-repository-url>
cd wallet-service
```

------------------------------------------------------------------------

### 2. Create Database

Login to PostgreSQL:

``` bash
psql -U postgres
```

or using any GUI like dbeaver

Create database:

``` sql
CREATE DATABASE wallet_db;
```

------------------------------------------------------------------------

### 3. Create Tables

Run the following SQL:

``` sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wallet_user_id ON wallets(user_id);
CREATE INDEX idx_transaction_user_id ON transactions(user_id);
```

------------------------------------------------------------------------

### 4. Seed Initial Data

``` sql
INSERT INTO users (name) VALUES ('Test User');

INSERT INTO wallets (user_id, balance)
VALUES (1, 100000);
```

------------------------------------------------------------------------

### 5. Configure Application

Create `config.json` in project root:

``` json
{
  "database": {
    "host": "localhost",
    "user": "postgres",
    "password": "your_password",
    "name": "wallet_db"
  }
}
```

------------------------------------------------------------------------

## Run Application

Install dependencies:

``` bash
go mod tidy
```

Run the server:

``` bash
go run cmd/api/main.go
```

Server will start at:

    http://localhost:8080

------------------------------------------------------------------------

## API Endpoints

### Health Check

    GET /ping

Response:

``` json
{
  "status": "ok",
  "database": "up"
}
```

If database is unreachable:

-   HTTP 503
-   Status: degraded

------------------------------------------------------------------------

### Withdraw

    POST /withdraw
    Content-Type: application/json

Request:

``` json
{
  "user_id": 1,
  "amount": 50000
}
```

Response:

``` json
{
  "message": "withdraw successful",
  "remaining_balance": 50000
}
```

Possible Errors:

-   400 → invalid input
-   400 → insufficient balance
-   404 → wallet not found
-   500 → internal error

------------------------------------------------------------------------

### Balance Inquiry

    GET /wallet?user_id=1

Response:

``` json
{
  "user_id": 1,
  "balance": 50000
}
```

------------------------------------------------------------------------

## Concurrency Handling

Withdraw operation uses:

-   Database transaction
-   `SELECT ... FOR UPDATE`
-   Row-level locking

This prevents race conditions and double-spending during concurrent
requests.

------------------------------------------------------------------------

## Author

Kevin Christian
