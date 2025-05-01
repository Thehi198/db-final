# Olin College Databases Final: Vector Databases

A simple Go-based vector database implementing approximate nearest neighbor search using the HNSW algorithm. It exposes CRUD operations via a library package and a REST API powered by [Fiber](https://github.com/gofiber/fiber).

---

## Table of Contents

- [Description](#description)
- [Installation & Bundling](#installation--bundling)
- [Running the Server](#running-the-server)
- [API Documentation](#api-documentation)
- [Library Functions](#library-functions)
- [Test Suite](#test-suite)

---

## Description

This project consists of two parts:

1. **Library (`vectordb/`)**: Implements:
   - `NewDatabase(dimension int) *Database` – create a new in-memory vector DB with HNSW index.
   - `(*Database).Insert(v Vector) (string, error)` – insert a new vector (checks dimension) and return its UUID.
   - `(*Database).Query(query []float64, k int, metadataFilter map[string]string) ([]Results, error)` – find k-NN using HNSW, optionally filter by metadata, then refine by exact Euclidean distance.
   - `(*Database).QueryByUUID(id string) (Vector, error)` – retrieve a vector by its UUID.
   - `(*Database).Update(id string, values []float64, metadata map[string]string) error` – update vector values and metadata.
   - `(*Database).Delete(id string) error` – remove a vector by its UUID.

2. **Server (`server.go`)**: A Fiber HTTP server exposing endpoints for:
   - **Create** a new DB instance (`POST /create`)
   - **Insert** vectors with metadata (`POST /insert`)
   - **Query** k nearest neighbors with metadata filter (`POST /query`)
   - **QueryByUUID** retrieve a vector by UUID (`POST /query_uuid`)
   - **Update** vector values and metadata (`PUT /update`)
   - **Delete** a vector by UUID (`DELETE /delete`)

---

## Installation & Bundling

```bash
# Clone project
cd /path/to/Project
# Ensure dependencies
go mod tidy

# Run tests
go test ./vectordb

# Build server binary
go build -o vectordb-server server.go

# Install FinnHub Library
go get github.com/Finnhub-Stock-API/finnhub-go/v2
```

---

## Running the Server

```bash
# Run directly to use the same terminal
go run server.go &

# To stop the server from running use
kill %1

# Or start built binary
./vectordb-server
```
The server listens on `http://localhost:3000`.

---

## API Documentation

We are using the Finnhub API for live financial data requests

```bash
# Flash your Finnhub private key onto bash
export FINNHUB_API_KEY=your_api_key_here

# Import Finnhub Go Library

import (
  finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

```

### 1. Create Database

- **Endpoint**: `POST /create`
- **Request JSON**:
  ```json
  { "dimension": 128, "metadata": { "key1": "value1" } }
  ```
- **Response**:
  ```json
  { "status": "ok", "dimension": 128, "metadata": { "key1": "value1" } }
  ```

### 2. Insert Vector

- **Endpoint**: `POST /insert`
- **Request JSON**:
  ```json
  { "values": [0.1, 0.2, ..., 0.128], "metadata": {"key1":"value1", ...} }
  ```
- **Response**:
  ```json
  { "uuid": "<generated-uuid>" }
  ```

### 3. Query k-NN

- **Endpoint**: `POST /query`
- **Request JSON**:
  ```json
  { "values": [0.1, 0.2, ..., 0.128], "k": 5, "metadata_filter": {"key1":"value1"} }
  ```
- **Response**:
  ```json
  [
    { "vector": { "uuid": "<uuid>", "values": [ ... ], "metadata": {"key1":"value1"} }, "distance": 0.123 },
    ...
  ]
  ```

### 4. QueryByUUID

- **Endpoint**: `POST /query_uuid`
- **Request JSON**:
  ```json
  { "uuid": "<existing-uuid>" }
  ```
- **Response**:
  ```json
  { "uuid": "<uuid>", "values": [ ... ], "metadata": {"key1":"value1"} }
  ```

### 5. Update Vector

- **Endpoint**: `PUT /update`
- **Request JSON**:
  ```json
  { "uuid": "<existing-uuid>", "values": [0.1, 0.2, ..., 0.128], "metadata": {"key1":"value1", ...} }
  ```
- **Response**:
  ```json
  { "status": "updated" }
  ```

### 6. Delete Vector

- **Endpoint**: `DELETE /delete`
- **Request JSON**:
  ```json
  { "uuid": "<existing-uuid>" }
  ```
- **Response**:
  ```json
  { "status": "deleted" }
  ```

---

## Library Functions

### Database Functions

```go 
// NewDatabase creates an in-memory vector DB with HNSW index for d-dimension.
func NewDatabase(d int) *Database

// Insert adds a new vector v to storage and index, returns its UUID.
func (db *Database) Insert(v Vector) (string, error)

// Query finds the k nearest neighbors to q, filters by metadataFilter if provided.
// Returns up to k results sorted by exact Euclidean distance.
func (db *Database) Query(q []float64, k int, metadataFilter map[string]string) ([]Results, error)

// QueryByUUID retrieves a stored vector by its UUID.
func (db *Database) QueryByUUID(id string) (Vector, error)

// Update modifies values and metadata of the vector with given UUID.
func (db *Database) Update(id string, values []float64, metadata map[string]string) error

// Delete removes the vector with given UUID.
func (db *Database) Delete(id string) error

// Vector represents a stored vector with UUID, values, and optional metadata.
type Vector struct {
    UUID     string
    Values   []float64
    Metadata map[string]string
}

// Results pairs a Vector with its distance to a query.
type Results struct {
    Vector   Vector
    Distance float64
}
```

### Vectorizing Functions

[Add section here to explain functions]

---

## Test Suite

All tests live in `vectordb/vectordb_test.go`:

1. **TestInsertDimensionMismatch**: ensures Insert rejects wrong dimensions.
2. **TestQueryDimensionMismatch**: ensures Query rejects wrong dimensions.
3. **TestQueryNearestNeighbor**: basic sanity check for 1-NN.
4. **TestQueryKGreaterThanLen**: `k > n` returns only `n` results.
5. **TestHNSWIndexLength**: HNSW index length matches inserted count.
6. **TestHNSWMatchesBruteForce**: HNSW top-k results appear in brute-force top-k.

Run:
```bash
go test ./vectordb -v
```

---
