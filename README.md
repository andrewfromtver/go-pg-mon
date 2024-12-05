# Postgres Query API

This API allows users to execute `SELECT` queries on a PostgreSQL database and returns the results in JSON format. It accepts input in JSON format and returns a JSON response with the query results.

## Endpoints

### `GET /swagger`

This endpoint serves the `Swagger UI` documentation for the API. It provides an interactive interface for users to explore and test the available API endpoints, including the ability to execute queries, view results, and understand the structure of requests and responses.

### `POST /db/custom-query`

This endpoint accepts a `POST` request to execute a `SELECT` query on the PostgreSQL database.

#### Request body

The request body must be a JSON object containing the following fields:

- **`dsn`** (string, required): The Data Source Name (DSN) for connecting to the PostgreSQL database. Example: `"postgres://[user]:[password]@[host]/[database]"`.
- **`query`** (string, required): The `SELECT` query to execute on the database.

##### Example - request body:
```json
{
  "dsn": "postgres://user:password@localhost:5432/dbname",
  "query": "SELECT id, name FROM users;"
}
```

##### Example - curl command:
```Bash
curl -X POST http://localhost:8080/db/custom-query \
  -H "Content-Type: application/json" \
  -d '{
    "dsn": "postgres://user:password@localhost:5432/dbname",
    "query": "SELECT id, name FROM users;"
  }'
```

### `POST /db/long-running-queries`

This endpoint executes a `SELECT` query on a `pg_stat_activity` to retrieve information about long-running queries. It allows you to dynamically specify parameters such as `query duration`, `database`, and `state` of the queries you want to monitor.

#### Request body

The request body must be a JSON object containing the following fields:

- **`dsn`** (string, required): The Data Source Name (DSN) for connecting to the PostgreSQL database. Example: `"postgres://[user]:[password]@[host]/[database]"`.
- **`query_duration`**  (string, required): The duration for filtering long-running queries. Only queries that have been running longer than this duration will be considered.
- **`pg_database`** (string, required): The name of the database.
- **`query_state`** (string, required): The state of the query (e.g., `"active"`, `"idle"`, `"idle in transaction"`). You can use `"*"` to include all states.

##### Example - request body:
```json
{
  "dsn": "postgres://user:password@localhost:5432/dbname",
  "query_duration": "5 second",
  "pg_database": "dbname",
  "query_state": "*"
}
```

##### Example - curl command:
```Bash
curl -X POST http://localhost:8080/db/long-running-queries \
  -H "Content-Type: application/json" \
  -d '{
    "dsn": "postgres://user:password@localhost:5432/dbname",
    "query_duration": "5 second",
    "pg_database": "dbname",
    "query_state": "*"
  }'
```

### `POST /db/largest-tables`

This endpoint executes a `SELECT` query on a `pg_statio_user_tables` to retrieve information about the largest tables.

#### Request body

The request body must be a JSON object containing the following fields:

- **`dsn`** (string, required): The Data Source Name (DSN) for connecting to the PostgreSQL database. Example: `"postgres://[user]:[password]@[host]/[database]"`.
- **`query_limit`**  (integer, required): The maximum number of results to return. Must be between `1` and `100`.

##### Example - request body:
```json
{
  "dsn": "postgres://user:password@localhost:5432/dbname",
  "query_limit": 10
}
```

##### Example - curl command:
```Bash
curl -X POST http://localhost:8080/db/largest-tables \
  -H "Content-Type: application/json" \
  -d '{
    "dsn": "postgres://user:password@localhost:5432/dbname",
    "query_limit": 10
  }'
```
