# Quirky Faulty Key-Value Store in GO

## Goal

This project is test assignment for ZOTA, and implements quirky key-value store in GO that randomly "breaks", operations have 30% chance to fault, simulating mysterious behavior.

---

## What It Does

Simple key-value store with three main operations:

- `Put(key, value)` — Stores data
- `Get(key)` — Retrieves data by key
- `Delete(key)` — Removes data by key

The data is persisted in MySQL.

---

## Sneaky faulty behaviour

Every operation has 30% chance to act incorrectly:

- `Get("cat")` might return wrong value by "fault".
- `Delete("apple")` might delete wrong record "key".
- `Put("key", "val")` might silently fail by "fault".

---

## Debugging Tool

The `Dump()` function prints actual internal state of the database to stdout, allowing to get overview of all current records.

---

## Extra Twist: Mutate Over Time

A background routine runs continuously, swapping values between two random keys every few seconds to simulate mutations and faulty behavior without `Put()` calls.

---

## How to Run

### Running Locally

1. Set required env vars by example:

    ```bash
    export APP_LOG_DIR="/path/to/log"
    export DB_HOST="localhost"
    export DB_USER="youruser"
    export DB_PASSWORD="yourpassword"
    export DB_NAME="yourdbname"
    ```

2. Run the program:

    ```bash
    go run main.go storage_tools.go
    ```

3. Program will:

    - Validate existanse of required vars
    - Initialize logging and DB
    - Insert example key-value pairs (`cat` → `meow`, `dog` → `woof`, `fish` → `blub`)
    - Retrieve sample value
    - Print the current DB state
    - Mutate data in the background indefinitely

---

### Running with Docker Compose

You can run the whole setup with Docker, which includes both GO app and MySQL server:

1. Make sure you have Docker installed.

2. Start services:

    ```bash
    docker-compose up
    ```

3. This will:

    - Run MySQL 8.0 container with database.
    - Run GO 1.24.3 application container with deffault environment variables.

4. As needed to execute tests instead of the app, you can adjust run command in `docker-compose.yml` to following:

    ```yaml
    command: >
      sh -c "go mod download && go test -v"
    ```


