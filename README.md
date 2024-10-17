# Multi Tenant API

Multi Tenant API is a backend service built with Go, using the Echo framework, designed to manage tenant data and integrate with RabbitMQ for message handling. The API allows you to create and delete tenants, consume messages from RabbitMQ queues, and provides logging with logrus.

## Features

- **Tenant Management**: Create and delete tenants.
- **RabbitMQ Integration**: Publish and consume messages with RabbitMQ.
- **Custom Logging**: Utilizes logrus for structured logging.
- **Performance Monitoring**: Middleware for tracking request performance.

## Requirements

- Go 1.17 or later
- PostgreSQL
- RabbitMQ

## Library & Documentation Library

- amqp091 (RabbitMq) (https://github.com/rabbitmq/amqp091-go)
- pqx (Database Postgress) (https://github.com/jackc/pgx)
- logrust (Logging) (https://github.com/sirupsen/logrus)
- lumberjack (Write Log) (https://github.com/natefinch/lumberjack)
- viper (Config) (https://github.com/spf13/viper)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/jatis_mobile_api.git
   cd jatis_mobile_api
   ```
2. Set up your Go environment::
   ```bash
   go mod tidy
   ```
3. Create a config.yaml file in the root directory with the following structure:
    ```yaml
    RabbitMQURL: "amqp://guest:guest@localhost:5672/"
    PostgresURL: "postgres://username:password@localhost:5432/database_name"
    PORT: 8080
    ```

## Running the Application

1. Make sure your PostgreSQL and RabbitMQ services are running.
2. Run the application:
    ```bash
    go run main.go
    ```
3. The server will start on the configured port (default: 8080).

## API Endpoints

### Create Tenant

- **POST** `/tenants`

**Request Body**:
```json
{
    "name": "Tenant Name"
}
```

### Response for Create Tenant

- **201 Created**: When the tenant is successfully created.
- **409 Conflict**: If the tenant already exists.

### Delete Tenant

- **DELETE** `/tenants/{id}`

**Response**:
- **200 OK**: When the tenant is successfully deleted.
- **400 Bad Request**: If the tenant ID is invalid.

### Consumer

- **GET** `/consumers`

**Headers**:
- `x-tenant-name: "Tenant Name"`

**Response**:
- **200 OK**: When the consumer is successfully started.

### Producer

- **POST** `/producers`

**Headers**:
- `x-tenant-name: "Tenant Name"`

**Request Body**:
```json
{
    "message": "Message for Tenant"
}
```

**Response**:
- **200 OK**: When the consumer is successfully started.

## Performance Monitoring
Middleware is included to log request performance metrics (duration, method, path).

## Logging
The application uses logrus for logging. Logs are written to app.log. The logging level can be adjusted as needed.

## Contributing
Contributions are welcome! Please submit a pull request or open an issue for any enhancements or bug fixes.

## License

### Instructions for Use
1. **Fill in the repository link**: Update the `git clone` line with the actual repository URL.
2. **Modify any sections** as necessary to fit your project's specifics, especially in the API Endpoints section if there are more or different endpoints.

This single markdown document provides a comprehensive overview of your project, setup instructions, API usage, testing, and contribution guidelines. Let me know if you need any further modifications or additions!
