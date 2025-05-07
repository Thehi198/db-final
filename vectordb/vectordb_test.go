// package vectordb

// import (
// 	"reflect"
// 	"sort"
// 	"testing"
// )

// // vectordb_test.go: detailed tests for vectordb operations and HNSW behavior.

// func TestInsertDimensionMismatch(t *testing.T) {
// 	db := NewDatabase(3)
// 	_, err := db.Insert(Vector{Values: []float64{1, 2}})
// 	t.Logf("Insert wrong-dimension error: %v", err)
// 	if err == nil {
// 		t.Errorf("expected error for wrong dimension, got nil")
// 	}
// }

// func TestQueryDimensionMismatch(t *testing.T) {
// 	db := NewDatabase(2)
// 	_, err := db.Query([]float64{1}, 1, nil)
// 	t.Logf("Query wrong-dimension error: %v", err)
// 	if err == nil {
// 		t.Errorf("expected error for wrong query dimension, got nil")
// 	}
// }

// func TestQueryNearestNeighbor(t *testing.T) {
// 	db := NewDatabase(2)
// 	v1 := Vector{Values: []float64{0, 0}}
// 	v2 := Vector{Values: []float64{1, 1}}
// 	t.Logf("Inserting vectors: %v and %v", v1.Values, v2.Values)
// 	_, _ = db.Insert(v1)
// 	_, _ = db.Insert(v2)
// 	t.Logf("Storage vectors: %v", db.Vectors)
// 	results, err := db.Query([]float64{1, 1}, 1, nil)
// 	t.Logf("Query([1,1], 1) returned: %+v", results)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	if !reflect.DeepEqual(results[0].Vector.Values, v2.Values) {
// 		t.Errorf("expected nearest vector %v, got %v", v2.Values, results[0].Vector.Values)
// 	}
// 	if results[0].Distance != 0 {
// 		t.Errorf("expected distance 0, got %v", results[0].Distance)
// 	}
// }

// func TestQueryKGreaterThanLen(t *testing.T) {
// 	db := NewDatabase(2)
// 	v := Vector{Values: []float64{0, 0}}
// 	t.Logf("Inserting single vector: %v", v.Values)
// 	_, _ = db.Insert(v)
// 	results, err := db.Query([]float64{0, 0}, 5, nil)
// 	t.Logf("Query(k>n) results length: %d", len(results))
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	if len(results) != 1 {
// 		t.Errorf("expected 1 result when k > len(vectors), got %d", len(results))
// 	}
// }

// // TestHNSWIndexLength verifies the HNSW index contains all inserted nodes.
// func TestHNSWIndexLength(t *testing.T) {
// 	db := NewDatabase(3)
// 	vs := []Vector{
// 		{Values: []float64{1, 2, 3}},
// 		{Values: []float64{4, 5, 6}},
// 		{Values: []float64{7, 8, 9}},
// 	}
// 	for _, v := range vs {
// 		if _, err := db.Insert(v); err != nil {
// 			t.Fatalf("unexpected insert error: %v", err)
// 		}
// 	}
// 	t.Logf("Inserted %d vectors; HNSW index length: %d", len(vs), db.index.Len())
// 	if got, want := db.index.Len(), len(vs); got != want {
// 		t.Errorf("expected index length %d, got %d", want, got)
// 	}
// }

// // TestHNSWMatchesBruteForce compares HNSW search against brute-force results.
// func TestHNSWMatchesBruteForce(t *testing.T) {
// 	db := NewDatabase(5)
// 	points := []Vector{
// 		{Values: []float64{0, 0, 0, 0, 0}},
// 		{Values: []float64{1, 1, 1, 1, 1}},
// 		{Values: []float64{2, 2, 2, 2, 2}},
// 		{Values: []float64{3, 3, 3, 3, 3}},
// 		{Values: []float64{4, 4, 4, 4, 4}},
// 	}
// 	for _, p := range points {
// 		if _, err := db.Insert(p); err != nil {
// 			t.Fatalf("insert error: %v", err)
// 		}
// 	}
// 	query := []float64{1.1, 1.1, 1.1, 1.1, 1.1}
// 	k := 3
// 	hnswRes, err := db.Query(query, k, nil)
// 	t.Logf("HNSW search results (vector, distance): %v", hnswRes)
// 	if err != nil {
// 		t.Fatalf("query error: %v", err)
// 	}
// 	type bf struct {
// 		idx  int
// 		dist float64
// 	}
// 	bfRes := make([]bf, len(db.Vectors))
// 	for i, v := range db.Vectors {
// 		bfRes[i] = bf{i, euclideanDistance(query, v.Values)}
// 	}
// 	sort.Slice(bfRes, func(i, j int) bool { return bfRes[i].dist < bfRes[j].dist })
// 	t.Logf("Brute-force sorted (idx, dist): %v", bfRes)
// 	// Ensure each HNSW result appears among the brute-force top-k, order may vary
// 	for _, hr := range hnswRes {
// 		found := false
// 		for j := 0; j < k; j++ {
// 			if reflect.DeepEqual(hr.Vector.Values, db.Vectors[bfRes[j].idx].Values) {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			t.Errorf("HNSW result %v not in brute-force top-%d", hr.Vector.Values, k)
// 		}
// 	}
// }

// // TestQueryWithMetadata ensures Query filters by metadata correctly.
// func TestQueryWithMetadata(t *testing.T) {
// 	db := NewDatabase(2)
// 	v1 := Vector{Values: []float64{0, 0}, Metadata: map[string]string{"type": "a"}}
// 	v2 := Vector{Values: []float64{1, 1}, Metadata: map[string]string{"type": "b"}}
// 	id1, _ := db.Insert(v1)
// 	_, _ = db.Insert(v2)
// 	// filter for type a
// 	results, err := db.Query([]float64{0, 0}, 2, map[string]string{"type": "a"})
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	if len(results) != 1 {
// 		t.Fatalf("expected 1 result for metadata filter, got %d", len(results))
// 	}
// 	if results[0].Vector.UUID != id1 {
// 		t.Errorf("expected UUID %s, got %s", id1, results[0].Vector.UUID)
// 	}
// 	if results[0].Vector.Metadata["type"] != "a" {
// 		t.Errorf("expected metadata type 'a', got %v", results[0].Vector.Metadata)
// 	}
// }

// // TestUpdateMetadata ensures Update correctly modifies vector values and metadata.
// func TestUpdateMetadata(t *testing.T) {
// 	db := NewDatabase(2)
// 	orig := Vector{Values: []float64{0, 0}, Metadata: map[string]string{"tag": "old"}}
// 	id, _ := db.Insert(orig)
// 	// update values and metadata
// 	err := db.Update(id, []float64{1, 1}, map[string]string{"tag": "new"})
// 	if err != nil {
// 		t.Fatalf("unexpected update error: %v", err)
// 	}
// 	// query by UUID to verify update
// 	updated, err := db.QueryByUUID(id)
// 	if err != nil {
// 		t.Fatalf("unexpected QueryByUUID error: %v", err)
// 	}
// 	if got, want := updated.Values, []float64{1, 1}; !reflect.DeepEqual(got, want) {
// 		t.Errorf("expected values %v, got %v", want, got)
// 	}
// 	if updated.Metadata["tag"] != "new" {
// 		t.Errorf("expected metadata tag 'new', got %v", updated.Metadata)
// 	}
// }

package vectordb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:3000"

func TestVectorDBIntegration(t *testing.T) {
	if !pingServer(t) {
		t.Fatal("Server is not running at localhost:3000")
	}

	initDatabase(t)
	insertRandomFinancialVectors(t)
	// fetchAndInsertFinancialVectors(t) // ⚠️ Uncomment to use real API
	printAllVectors(t)

	t.Run("QueryExistingData", TestQueryExistingData)
	t.Run("QueryByMetadata", TestQueryByMetadata)
	t.Run("QueryByUUID", TestQueryByUUID)
}

func pingServer(t *testing.T) bool {
	resp, err := http.Get(baseURL + "/list")
	if err != nil {
		t.Logf("Ping failed: %v", err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func initDatabase(t *testing.T) {
	payload := map[string]interface{}{
		"dimension": 4,
	}
	resp, err := sendPOST(baseURL+"/create", payload)
	if err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}
	t.Logf("DB init response: %s", resp)
}

func insertRandomFinancialVectors(t *testing.T) {
	tickers := []string{"AAPL", "MSFT", "GOOG", "TSLA", "NVDA", "META", "AMZN"}
	rand.Seed(time.Now().UnixNano())
	for _, ticker := range tickers {
		values := []float64{
			rand.Float64()*1000 + 100,
			rand.Float64()*1000 + 100,
			rand.Float64()*1000 + 100,
			rand.Float64()*1000 + 100,
		}
		metadata := map[string]string{"ticker": ticker}
		payload := map[string]interface{}{
			"values":   values,
			"metadata": metadata,
		}
		result, err := sendPOST(baseURL+"/insert", payload)
		if err != nil {
			t.Fatalf("insert error for %s: %v", ticker, err)
		}
		t.Logf("Inserted %s: %s", ticker, result)
	}
}

func TestQueryExistingData(t *testing.T) {
	query := map[string]interface{}{
		"values": []float64{688.0, 692.3, 685.4, 690.1},
		"k":      2,
	}
	jsonData, _ := json.Marshal(query)
	resp, err := http.Post(baseURL+"/query", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	var results []map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&results); err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("decode error: %v\nRaw: %s", err, string(body))
	}
	t.Logf("Query result: %+v", results)
}

func TestQueryByMetadata(t *testing.T) {
	query := map[string]interface{}{
		"values": []float64{143.2, 144.1, 142.7, 143.8},
		"k":      10,
		"metadata_filter": map[string]string{
			"ticker": "AAPL",
		},
	}
	jsonData, _ := json.Marshal(query)
	resp, err := http.Post(baseURL+"/query", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	var results []map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&results); err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("decode error: %v\nRaw: %s", err, string(body))
	}
	t.Logf("Metadata query result: %+v", results)
}

func TestQueryByUUID(t *testing.T) {
	uuid := getUUIDForTicker(t, "AAPL")
	if uuid == "" {
		t.Skip("No vector found for AAPL")
	}
	query := map[string]interface{}{
		"uuid": uuid,
	}
	jsonData, _ := json.Marshal(query)
	resp, err := http.Post(baseURL+"/query_uuid", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("decode error: %v\nRaw: %s", err, string(body))
	}
	t.Logf("UUID query result: %+v", result)
}

func getUUIDForTicker(t *testing.T, ticker string) string {
	resp, err := http.Get(baseURL + "/list")
	if err != nil {
		t.Errorf("failed to fetch vectors: %v", err)
		return ""
	}
	defer resp.Body.Close()

	var vectors []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&vectors); err != nil {
		t.Errorf("decode error: %v", err)
		return ""
	}

	for _, v := range vectors {
		if meta, ok := v["metadata"].(map[string]interface{}); ok {
			if tickerVal, ok := meta["ticker"].(string); ok && tickerVal == ticker {
				if uuid, ok := v["uuid"].(string); ok {
					return uuid
				}
			}
		}
	}
	return ""
}

func printAllVectors(t *testing.T) {
	resp, err := http.Get(baseURL + "/list")
	if err != nil {
		t.Errorf("failed to fetch all vectors: %v", err)
		return
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		t.Errorf("decode error: %v", err)
		return
	}

	fmt.Println("\n--- All Vectors in DB ---")
	for i, v := range result {
		fmt.Printf("[%d] UUID: %s, Metadata: %v, Values: %v\n", i, v["uuid"], v["metadata"], v["values"])
	}
}

func sendPOST(url string, payload map[string]interface{}) (string, error) {
	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var generic interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &generic); err != nil {
		return "", fmt.Errorf("decode error: %v\nRaw body: %s", err, string(body))
	}

	out, _ := json.MarshalIndent(generic, "", "  ")
	return string(out), nil
}
