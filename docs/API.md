# API Documentation

This document provides comprehensive documentation for the Auth Service API. The API follows RESTful principles and uses JSON for request and response bodies.

## Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
- [Error Handling](#error-handling)
- [Health Checks](#health-checks)
  - [Health Check](#health-check)
  - [Liveness Probe](#liveness-probe)
  - [Readiness Probe](#readiness-probe)
- [API v1](#api-v1)
  - [Ping](#ping)
- [Rate Limiting](#rate-limiting)
- [Response Format](#response-format)
- [Status Codes](#status-codes)
- [Versioning](#versioning)

## Base URL

All API endpoints are relative to the base URL of the service:

```
http://localhost:8080
```

In production, replace `localhost:8080` with your actual domain and port.

## Authentication

> **Note**: Authentication endpoints will be implemented in future iterations.

## Error Handling

Errors are returned in JSON format with the following structure:

```json
{
  "error": {
    "code": "error_code",
    "message": "Human-readable error message",
    "details": {
      "field1": "Error detail 1",
      "field2": "Error detail 2"
    }
  },
  "timestamp": "2025-01-01T12:00:00Z"
}
```

## Health Checks

### Health Check

Get the health status of the service and its dependencies.

**Endpoint**: `GET /health`

**Response**:

```json
{
  "status": "healthy",
  "timestamp": "2025-01-01T12:00:00Z",
  "checks": [
    {
      "name": "server",
      "status": "healthy",
      "message": "Server is running",
      "timestamp": "2025-01-01T12:00:00Z"
    },
    {
      "name": "database",
      "status": "healthy",
      "message": "Database is connected",
      "timestamp": "2025-01-01T12:00:00Z"
    },
    {
      "name": "logger",
      "status": "healthy",
      "message": "Logger is operational",
      "timestamp": "2025-01-01T12:00:00Z"
    }
  ]
}
```

**Status Codes**:
- `200 OK`: Service is healthy or degraded
- `503 Service Unavailable`: Service is unhealthy

### Liveness Probe

Simple endpoint to check if the service is running.

**Endpoint**: `GET /health/live`

**Response**:

```json
{
  "status": "ok",
  "timestamp": "2025-01-01T12:00:00Z"
}
```

**Status Codes**:
- `200 OK`: Service is running

### Readiness Probe

Check if the service is ready to accept requests.

**Endpoint**: `GET /health/ready`

**Response**:

```json
{
  "status": "ready",
  "timestamp": "2025-01-01T12:00:00Z"
}
```

**Status Codes**:
- `200 OK`: Service is ready to accept requests
- `503 Service Unavailable`: Service is not ready

## API v1

### Ping

Simple endpoint to test the API connectivity.

**Endpoint**: `GET /api/v1/ping`

**Response**:

```json
{
  "message": "pong",
  "time": "2025-01-01T12:00:00Z"
}
```

**Status Codes**:
- `200 OK`: Request successful

## Rate Limiting

> **Note**: Rate limiting will be implemented in a future iteration.

## Response Format

All successful API responses follow this format:

```json
{
  "data": {},
  "meta": {
    "timestamp": "2025-01-01T12:00:00Z",
    "version": "1.0.0"
  }
}
```

- `data`: Contains the response payload
- `meta`: Contains metadata about the response

## Status Codes

The API uses the following status codes:

| Code | Status                  | Description                                      |
|------|-------------------------|--------------------------------------------------|
| 200  | OK                      | Request was successful                           |
| 201  | Created                 | Resource created successfully                    |
| 204  | No Content              | Request successful, no content to return         |
| 400  | Bad Request             | Invalid request format or parameters             |
| 401  | Unauthorized            | Authentication required                          |
| 403  | Forbidden               | Insufficient permissions                         |
| 404  | Not Found               | Resource not found                               |
| 405  | Method Not Allowed      | HTTP method not allowed for the requested route  |
| 409  | Conflict                | Resource conflict                                |
| 422  | Unprocessable Entity    | Validation failed                                |
| 429  | Too Many Requests       | Rate limit exceeded                              |
| 500  | Internal Server Error   | Server error                                     |
| 502  | Bad Gateway             | Upstream service unavailable                     |
| 503  | Service Unavailable     | Service temporarily unavailable                  |
| 504  | Gateway Timeout         | Upstream service timeout                         |

## Versioning

The API uses URL versioning (e.g., `/api/v1/...`). Breaking changes will result in a new version number.

## Pagination

> **Note**: Pagination will be implemented in a future iteration for list endpoints.

## Filtering and Sorting

> **Note**: Filtering and sorting will be implemented in a future iteration.

## WebSocket API

> **Note**: WebSocket support will be added in a future iteration.

## Webhooks

> **Note**: Webhook support will be added in a future iteration.
