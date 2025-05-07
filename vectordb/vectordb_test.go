package vectordb

import (
	"reflect"
	"sort"
	"testing"
)

// vectordb_test.go: detailed tests for vectordb operations and HNSW behavior.

func TestInsertDimensionMismatch(t *testing.T) {
	db := NewDatabase(3)
	_, err := db.Insert(Vector{Values: []float64{1, 2}})
	t.Logf("Insert wrong-dimension error: %v", err)
	if err == nil {
		t.Errorf("expected error for wrong dimension, got nil")
	}
}

func TestQueryDimensionMismatch(t *testing.T) {
	db := NewDatabase(2)
	_, err := db.Query([]float64{1}, 1, nil)
	t.Logf("Query wrong-dimension error: %v", err)
	if err == nil {
		t.Errorf("expected error for wrong query dimension, got nil")
	}
}

func TestQueryNearestNeighbor(t *testing.T) {
	db := NewDatabase(2)
	v1 := Vector{Values: []float64{0, 0}}
	v2 := Vector{Values: []float64{1, 1}}
	t.Logf("Inserting vectors: %v and %v", v1.Values, v2.Values)
	_, _ = db.Insert(v1)
	_, _ = db.Insert(v2)
	t.Logf("Storage vectors: %v", db.Vectors)
	results, err := db.Query([]float64{1, 1}, 1, nil)
	t.Logf("Query([1,1], 1) returned: %+v", results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(results[0].Vector.Values, v2.Values) {
		t.Errorf("expected nearest vector %v, got %v", v2.Values, results[0].Vector.Values)
	}
	if results[0].Distance != 0 {
		t.Errorf("expected distance 0, got %v", results[0].Distance)
	}
}

func TestQueryKGreaterThanLen(t *testing.T) {
	db := NewDatabase(2)
	v := Vector{Values: []float64{0, 0}}
	t.Logf("Inserting single vector: %v", v.Values)
	_, _ = db.Insert(v)
	results, err := db.Query([]float64{0, 0}, 5, nil)
	t.Logf("Query(k>n) results length: %d", len(results))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result when k > len(vectors), got %d", len(results))
	}
}

// TestHNSWIndexLength verifies the HNSW index contains all inserted nodes.
func TestHNSWIndexLength(t *testing.T) {
	db := NewDatabase(3)
	vs := []Vector{
		{Values: []float64{1, 2, 3}},
		{Values: []float64{4, 5, 6}},
		{Values: []float64{7, 8, 9}},
	}
	for _, v := range vs {
		if _, err := db.Insert(v); err != nil {
			t.Fatalf("unexpected insert error: %v", err)
		}
	}
	t.Logf("Inserted %d vectors; HNSW index length: %d", len(vs), db.index.Len())
	if got, want := db.index.Len(), len(vs); got != want {
		t.Errorf("expected index length %d, got %d", want, got)
	}
}

// TestHNSWMatchesBruteForce compares HNSW search against brute-force results.
func TestHNSWMatchesBruteForce(t *testing.T) {
	db := NewDatabase(5)
	points := []Vector{
		{Values: []float64{0, 0, 0, 0, 0}},
		{Values: []float64{1, 1, 1, 1, 1}},
		{Values: []float64{2, 2, 2, 2, 2}},
		{Values: []float64{3, 3, 3, 3, 3}},
		{Values: []float64{4, 4, 4, 4, 4}},
	}
	for _, p := range points {
		if _, err := db.Insert(p); err != nil {
			t.Fatalf("insert error: %v", err)
		}
	}
	query := []float64{1.1, 1.1, 1.1, 1.1, 1.1}
	k := 3
	hnswRes, err := db.Query(query, k, nil)
	t.Logf("HNSW search results (vector, distance): %v", hnswRes)
	if err != nil {
		t.Fatalf("query error: %v", err)
	}
	type bf struct {
		idx  int
		dist float64
	}
	bfRes := make([]bf, len(db.Vectors))
	for i, v := range db.Vectors {
		bfRes[i] = bf{i, euclideanDistance(query, v.Values)}
	}
	sort.Slice(bfRes, func(i, j int) bool { return bfRes[i].dist < bfRes[j].dist })
	t.Logf("Brute-force sorted (idx, dist): %v", bfRes)
	// Ensure each HNSW result appears among the brute-force top-k, order may vary
	for _, hr := range hnswRes {
		found := false
		for j := 0; j < k; j++ {
			if reflect.DeepEqual(hr.Vector.Values, db.Vectors[bfRes[j].idx].Values) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("HNSW result %v not in brute-force top-%d", hr.Vector.Values, k)
		}
	}
}

// TestQueryWithMetadata ensures Query filters by metadata correctly.
func TestQueryWithMetadata(t *testing.T) {
	db := NewDatabase(2)
	v1 := Vector{Values: []float64{0, 0}, Metadata: map[string]string{"type": "a"}}
	v2 := Vector{Values: []float64{1, 1}, Metadata: map[string]string{"type": "b"}}
	id1, _ := db.Insert(v1)
	_, _ = db.Insert(v2)
	// filter for type a
	results, err := db.Query([]float64{0, 0}, 2, map[string]string{"type": "a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result for metadata filter, got %d", len(results))
	}
	if results[0].Vector.UUID != id1 {
		t.Errorf("expected UUID %s, got %s", id1, results[0].Vector.UUID)
	}
	if results[0].Vector.Metadata["type"] != "a" {
		t.Errorf("expected metadata type 'a', got %v", results[0].Vector.Metadata)
	}
}

// TestUpdateMetadata ensures Update correctly modifies vector values and metadata.
func TestUpdateMetadata(t *testing.T) {
	db := NewDatabase(2)
	orig := Vector{Values: []float64{0, 0}, Metadata: map[string]string{"tag": "old"}}
	id, _ := db.Insert(orig)
	// update values and metadata
	err := db.Update(id, []float64{1, 1}, map[string]string{"tag": "new"})
	if err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}
	// query by UUID to verify update
	updated, err := db.QueryByUUID(id)
	if err != nil {
		t.Fatalf("unexpected QueryByUUID error: %v", err)
	}
	if got, want := updated.Values, []float64{1, 1}; !reflect.DeepEqual(got, want) {
		t.Errorf("expected values %v, got %v", want, got)
	}
	if updated.Metadata["tag"] != "new" {
		t.Errorf("expected metadata tag 'new', got %v", updated.Metadata)
	}
}
