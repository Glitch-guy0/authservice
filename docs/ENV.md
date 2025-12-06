# Environment Variables

This document outlines all the environment variables used in the Auth Service. These variables can be set either in a `.env` file in the project root or as system environment variables.

## Table of Contents

- [Application Environment](#application-environment)
- [Server Configuration](#server-configuration)
- [Logging](#logging)
- [Database](#database)

## Application Environment

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | `development` | Application environment. Can be `development`, `staging`, or `production`. |

## Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Port on which the server will listen for incoming connections. |
| `SERVER_READ_TIMEOUT` | `15` | Maximum duration in seconds for reading the entire request, including the body. |
| `SERVER_WRITE_TIMEOUT` | `15` | Maximum duration in seconds before timing out writes of the response. |
| `SERVER_IDLE_TIMEOUT` | `60` | Maximum amount of time in seconds to wait for the next request when keep-alives are enabled. |
| `SERVER_DEBUG` | `false` | Enable debug mode for the server. |

## Logging

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Logging level. Can be `debug`, `info`, `warn`, `error`, `fatal`, or `panic`. |
| `LOG_FORMAT` | `json` | Log format. Can be `json` or `text`. |
| `LOG_FILE_ENABLED` | `false` | Enable or disable file logging. |
| `LOG_FILENAME` | `app.log` | Name of the log file. |
| `LOG_MAX_SIZE` | `100` | Maximum size in MB of the log file before it gets rotated. |
| `LOG_MAX_AGE` | `30` | Maximum number of days to retain old log files. |

## Database

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_DRIVER` | `postgres` | Database driver to use. Supported values: `postgres`, `mysql`, `sqlite`. |
| `DB_DSN` | `host=localhost port=5432 user=postgres password=postgres dbname=auth_service sslmode=disable` | Database connection string. |
| `DB_MAX_OPEN_CONNS` | `25` | Maximum number of open connections to the database. |
| `DB_MAX_IDLE_CONNS` | `5` | Maximum number of idle connections to the database. |
| `DB_CONN_MAX_LIFETIME` | `5` | Maximum amount of time in minutes a connection may be reused. |
| `DB_LOG_QUERIES` | `false` | Enable or disable logging of database queries. |

## Using Environment Variables

### In Development

1. Create a `.env` file in the project root:
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   # Edit the config.yaml file as needed
   ```

2. Or set environment variables in your shell:
   ```bash
   export APP_ENV=development
   export SERVER_PORT=3000
   export DB_DSN="your_connection_string_here"
   ```

### In Production

For production, it's recommended to use system environment variables or a secrets management system:

```bash
# Example using environment variables in a production environment
export APP_ENV=production
export SERVER_PORT=80
export DB_DSN="host=db.example.com port=5432 user=prod_user password=secure_password dbname=auth_prod sslmode=require"
```

### In Docker

When using Docker, you can pass environment variables in your `docker-compose.yml`:

```yaml
services:
  auth-service:
    image: your-auth-service:latest
    environment:
      - APP_ENV=production
      - SERVER_PORT=8080
      - DB_DSN=postgres://user:password@db:5432/auth_service?sslmode=disable
    depends_on:
      - db
```

## Security Considerations

- Never commit sensitive information like database credentials in version control.
- Use a `.gitignore` file to exclude `.env` files from being committed.
- In production, use a secrets management solution or environment variables provided by your hosting platform.
- Rotate database credentials and API keys regularly.
- Set appropriate file permissions for configuration files containing sensitive information.

## Troubleshooting

- If environment variables are not being loaded, ensure they are set before starting the application.
- Check for typos in variable names.
- Remember that environment variable names are case-sensitive.
- For file logging issues, verify that the application has write permissions to the log directory.
