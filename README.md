# Postgres Query API

This API allows users to execute `SELECT` queries on a PostgreSQL database and returns the results in JSON format. It accepts input in JSON format and returns a JSON response with the query results.

## Endpoints

### `POST /query`

This endpoint accepts a `POST` request to execute a `SELECT` query on the PostgreSQL database.

#### Request body

The request body must be a JSON object containing the following fields:

- **`dsn`** (string, required): The Data Source Name (DSN) for connecting to the PostgreSQL database. Example: `"postgres://[user]:[password]@[host]/[database]"`.
- **`query`** (string, required): The `SELECT` query to execute on the database.
- **`output`** (string, required): The name of the output file where the results will be saved. This can be any string, though this is not used in the current implementation but expected in the request.

##### Example - request body:
```json
{
  "dsn": "postgres://user:password@localhost:5432/dbname",
  "query": "SELECT id, name FROM users;",
  "output": "output.json"
}
```

##### Example - curl command:
```Bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "dsn": "postgres://user:password@localhost:5432/dbname",
    "query": "SELECT id, name FROM users;",
    "output": "output.json"
  }'
```
