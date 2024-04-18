# Payments Service

This project is a payments service implemented in Go, using PostgreSQL as the database, all containerized using Docker.

## Overview

The application provides a RESTful API to manage financial transactions, including:
- Creating accounts with initial balances.
- Transferring funds between accounts.
- Retrieving account balances.

## Assumptions

This project is built with several key assumptions that influence its design and functionality:

- **Single-Currency Transactions**: The application is designed to handle financial transactions in a single currency. This simplifies the handling of monetary values and avoids the complexities related to currency conversion and fluctuation.

- **PostgreSQL Database**: The system assumes the use of PostgreSQL for data storage. PostgreSQL's features such as ACID compliance, transaction management, and robustness are integral to the application's data handling strategies.

- **SERIALIZABLE Isolation Level**: All transactions within the database are handled at the `SERIALIZABLE` isolation level. This is the highest level of isolation and is assumed to prevent lost updates, dirty reads, and other transaction anomalies, at the cost of potential performance overhead due to increased locking and conflict resolution.

- **Security Practices**: Basic security practices are assumed to be in place, including secure handling of user data and credentials. However, specific details like OAuth implementation or advanced security measures may need to be tailored according to user needs.


## Getting Started

### Prerequisites

- Docker and Docker Compose must be installed on your machine.

### Setting Up and Running

The application and its database can be started using Docker Compose:

```bash
make up
```
This command builds the Docker images if necessary and starts the containers specified in the docker-compose.yml file. The payments service depends on the db service, ensuring the database is ready before the application starts.

### Stopping the Services

To stop and remove the containers:

```bash
make down
```

### Environment Variables

The following environment variables are crucial for connecting the application to PostgreSQL:

- **DB_HOST**: Hostname of the PostgreSQL server.
- **DB_PORT**: Port on which PostgreSQL is accessible.
- **DB_USER**: Username for the PostgreSQL database.
- **DB_PASSWORD**: Password for the PostgreSQL database.
- **DB_NAME**: Database name used by the application.

These variables are set in the docker-compose.yml for Docker-based setups and should be configured accordingly for local environments in your development setup or IDE.

## API Documentation

- **Create Account**: POST /accounts with JSON body {"account_id": "123", "balance": 100.00}
- **Get Account**: GET /accounts/123
- **Transfer Funds**: POST /transactions with JSON body {"source_account_id": "123", "destination_account_id": "234", "amount": 50.00}


## Local Development
For developers working directly within the payments directory:

### Building Locally
Navigate to the project directory and run:

```bash
make local-build
```
This command builds the Go application within the payments subdirectory.

Running Locally
To run the built application:

```bash
make local-run
```
This command executes the binary compiled from the payments source code.

### Testing
Running Tests Locally
To execute tests in the local environment:

```bash
make local-test
```