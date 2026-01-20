# go-mapper-showcase

A REST API demonstrating the [tariklabs/mapper](https://github.com/tariklabs/mapper) library features

## Mapper Features Demonstrated

| Feature                  | Description                                               |
| ------------------------ | --------------------------------------------------------- |
| `Map(dst, src)`          | Basic struct-to-struct mapping                            |
| `MapWithOptions()`       | Mapping with custom configuration                         |
| `map` tag                | Field aliasing (`lat` → `Latitude`)                       |
| `mapconv` tag            | Type conversion (`string` → `float64`)                    |
| `WithStrictMode()`       | Ensures all destination fields have matching source       |
| `WithIgnoreZeroSource()` | Skips zero-value fields (useful for partial updates)      |
| `WithMaxDepth()`         | Sets maximum nesting depth for recursive mapping          |
| `MappingError`           | Structured error handling with field path info            |
| Nested structs           | Recursive mapping of complex objects                      |

## Project Structure

```
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── app/
│   │   └── app.go           # Routes and dependency wiring
│   ├── core/
│   │   ├── domain/          # Domain entities
│   │   ├── handler/         # HTTP handlers (DTOs + mapper usage)
│   │   ├── repository/      # Data access (DTOs + mapper usage)
│   │   └── service/         # Business logic
│   └── engine/
│       └── engine.go        # HTTP helper utilities
```

## Running

```bash
go run ./cmd/main.go
```

Server starts at `http://localhost:8080`

## Endpoints

### GET /users/:id

Fetches a user from JSONPlaceholder and maps the response using nested struct mapping with field aliasing and type conversion.

```bash
curl http://localhost:8080/users/1
```

### POST /users

Creates a user with mapper handling request-to-domain conversion.

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "username": "johnd",
    "email": "john@example.com",
    "phone": "1-555-123-4567",
    "website": "johndoe.com",
    "address": {
      "street": "123 Main St",
      "suite": "Apt 4B",
      "city": "New York",
      "zipcode": "10001",
      "geo": { "lat": "40.7128", "lng": "-74.0060" }
    },
    "company": {
      "name": "Doe Industries",
      "catch_phrase": "Innovating the future",
      "bs": "synergize scalable solutions"
    }
  }'
```

## Key Mapping Examples

**String to Float64 conversion:**

```go
type GeoDTO struct {
    Lat string `map:"Latitude" mapconv:"float64"`
    Lng string `map:"Longitude" mapconv:"float64"`
}
```

**Field aliasing:**

```go
type CompanyDTO struct {
    CatchPhrase string `map:"CatchPhrase"`
    BS          string `map:"BS"`
}
```

**Mapping with options:**

```go
mapper.MapWithOptions(&response, user,
    mapper.WithStrictMode(),
    mapper.WithMaxDepth(10),
)
```
