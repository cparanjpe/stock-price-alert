Here's a comprehensive `README.md` file that includes all the necessary details to get this Go project up and running:

---

# BTC/USDT Price Alert System

This project is a Go-based price alert system for the BTC/USDT trading pair using Binance WebSockets. The system allows users to set alerts based on specific conditions, such as price movements and technical indicators (RSI/MACD). Alerts are stored in MySQL for persistence and cached in Redis for optimization.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Project Setup](#project-setup)
  - [1. Clone the Repository](#1-clone-the-repository)
  - [2. Install Go Dependencies](#2-install-go-dependencies)
  - [3. Setup MySQL](#3-setup-mysql)
  - [4. Setup Redis](#4-setup-redis)
  - [5. Run the Application](#5-run-the-application)
- [MySQL Schema](#mysql-schema)
- [Docker Commands](#docker-commands)
- [Project Structure](#project-structure)

## Prerequisites

Ensure that you have the following installed on your machine:

- [Go 1.19+](https://golang.org/doc/install)
- [MySQL 8.0+](https://dev.mysql.com/downloads/mysql/)
- [Docker](https://www.docker.com/get-started)

## Project Setup

### 1. Clone the Repository

```bash
git clone https://github.com/cparanjpe/btcusdt-alert.git
cd btcusdt-alert
```

### 2. Install Go Dependencies

Ensure all the required Go modules are downloaded:

```bash
go mod tidy
```

### 3. Setup MySQL

1. Start MySQL server (if not already running).

2. Create a database for the project:

    ```sql
    CREATE DATABASE yourdatabase;
    ```

3. Create the `alerts` table using the following SQL DDL statement:

    ```sql
    CREATE TABLE alerts (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id VARCHAR(255) NOT NULL,
        value FLOAT NOT NULL,
        direction ENUM('up', 'down') NOT NULL,
        indicator ENUM('rsi', 'macd') NOT NULL,
        alerted BOOLEAN NOT NULL DEFAULT FALSE,
        UNIQUE KEY unique_alert (user_id, value, direction, indicator)
    );
    ```

4. Update the MySQL connection details in the `database/db.go` file:

    ```go
    db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/yourdatabase")
    ```

### 4. Setup Redis

You can run Redis using Docker. Follow the steps below:

1. Pull the Redis Docker image:

    ```bash
    docker pull redis
    ```

2. Run Redis using Docker:

    ```bash
    docker run --name redis-server -p 6379:6379 -d redis
    ```

### 5. Run the Application

With everything set up, you can run the Go application:

```bash
go run main.go
```

## MySQL Schema

Here’s the schema for the `alerts` table used in this project:

```sql
CREATE TABLE alerts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    value FLOAT NOT NULL,
    direction ENUM('up', 'down') NOT NULL,
    indicator ENUM('rsi', 'macd') NOT NULL,
    alerted BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE KEY unique_alert (user_id, value, direction, indicator)
);
```

## Docker Commands

To set up and run Redis using Docker:

1. Pull the Redis image:

    ```bash
    docker pull redis
    ```

2. Run the Redis container:

    ```bash
    docker run --name redis-server -p 6379:6379 -d redis
    ```

3. To stop the Redis container:

    ```bash
    docker stop redis-server
    ```

4. To remove the Redis container:

    ```bash
    docker rm redis-server
    ```

## Project Structure

- `main.go`: Entry point of the application.
- `alert/alert.go`: Contains logic for managing alerts, including adding, processing, and fetching alerts.
- `database/db.go`: Handles MySQL database connection initialization.
- `notifier/notifier.go`: Contains logic for notifying users when an alert is triggered.

---

By following these instructions, you should be able to set up and run the BTC/USDT Price Alert System on your local machine. If you encounter any issues, please consult the documentation or reach out for support.
