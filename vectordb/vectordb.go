package vectordb

import (
	"fmt"
	"math"
	"sort"

	"github.com/coder/hnsw"
	"github.com/google/uuid"
)

type Vector struct {
	UUID     string            `json:"uuid"`
	Values   []float64         `json:"values"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type Results struct {
	Vector   Vector  `json:"vector"`
	Distance float64 `json:"distance"`
}

type Database struct {
	Dimension int
	Metadata  map[string]string
	Vectors   []Vector
	index     *hnsw.Graph[int]
}

// NewDatabase creates a new vector database with the given dimension.
func NewDatabase(dimension int) *Database {
	idx := hnsw.NewGraph[int]()
	// set Euclidean distance for HNSW
	idx.Distance = func(a, b hnsw.Vector) float32 {
		var sum float32
		for i := range a {
			diff := a[i] - b[i]
			sum += diff * diff
		}
		return float32(math.Sqrt(float64(sum)))
	}
	return &Database{
		Dimension: dimension,
		Metadata:  make(map[string]string),
		Vectors:   make([]Vector, 0),
		index:     idx,
	}
}

// Insert adds a new vector to the database.
func (db *Database) Insert(v Vector) (string, error) {
	if len(v.Values) != db.Dimension {
		return "", fmt.Errorf("vector dimension %d does not match database dimension %d", len(v.Values), db.Dimension)
	}
	id := uuid.New().String()
	v.UUID = id
	// append to storage
	db.Vectors = append(db.Vectors, v)
	// add to HNSW index
	key := len(db.Vectors) - 1
	db.index.Add(hnsw.MakeNode(key, float64ToFloat32Slice(v.Values)))
	return id, nil
}

// Query returns the k nearest neighbors to the query vector using HNSW index with optional metadata filtering.
func (db *Database) Query(query []float64, k int, metadataFilter map[string]string) ([]Results, error) {
	if len(query) != db.Dimension {
		return nil, fmt.Errorf("query dimension %d does not match database dimension %d", len(query), db.Dimension)
	}
	// approximate search via HNSW: retrieve all nodes for full candidate set
	f32Query := float64ToFloat32Slice(query)
	nodes := db.index.Search(f32Query, len(db.Vectors))
	// collect candidates with metadata filter
	var candidates []Results
	for _, node := range nodes {
		v := db.Vectors[node.Key]
		if metadataFilter != nil && len(metadataFilter) > 0 {
			match := true
			for mk, mv := range metadataFilter {
				if v.Metadata == nil || v.Metadata[mk] != mv {
					match = false
					break
				}
			}
			if !match {
				continue
			}
		}
		candidates = append(candidates, Results{
			Vector:   v,
			Distance: euclideanDistance(query, v.Values),
		})
	}
	// sort by actual distance
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Distance < candidates[j].Distance
	})
	// trim to k results
	if len(candidates) > k {
		candidates = candidates[:k]
	}
	return candidates, nil
}

// euclideanDistance computes the Euclidean distance between two vectors of same length.
func euclideanDistance(a, b []float64) float64 {
	var sum float64
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}

// float64ToFloat32Slice converts []float64 to []float32 for HNSW index.
func float64ToFloat32Slice(in []float64) []float32 {
	out := make([]float32, len(in))
	for i, v := range in {
		out[i] = float32(v)
	}
	return out
}

// QueryByUUID returns a vector by its UUID.
func (db *Database) QueryByUUID(id string) (Vector, error) {
	for _, v := range db.Vectors {
		if v.UUID == id {
			return v, nil
		}
	}
	return Vector{}, fmt.Errorf("uuid %s not found", id)
}

// Delete removes a vector by its UUID.
func (db *Database) Delete(id string) error {
	idx := -1
	for i, v := range db.Vectors {
		if v.UUID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("uuid %s not found", id)
	}
	// remove from slice
	db.Vectors = append(db.Vectors[:idx], db.Vectors[idx+1:]...)
	// rebuild HNSW index
	newIdx := hnsw.NewGraph[int]()
	newIdx.Distance = db.index.Distance
	for i, v := range db.Vectors {
		newIdx.Add(hnsw.MakeNode(i, float64ToFloat32Slice(v.Values)))
	}
	db.index = newIdx
	return nil
}

// Update modifies a vector's values and metadata by its UUID.
func (db *Database) Update(id string, newValues []float64, metadata map[string]string) error {
	if len(newValues) != db.Dimension {
		return fmt.Errorf("vector dimension %d does not match database dimension %d", len(newValues), db.Dimension)
	}
	idx := -1
	for i, v := range db.Vectors {
		if v.UUID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("uuid %s not found", id)
	}
	db.Vectors[idx].Values = newValues
	// update metadata
	db.Vectors[idx].Metadata = metadata
	// rebuild HNSW index
	newIdx := hnsw.NewGraph[int]()
	newIdx.Distance = db.index.Distance
	for i, v := range db.Vectors {
		newIdx.Add(hnsw.MakeNode(i, float64ToFloat32Slice(v.Values)))
	}
	db.index = newIdx
	return nil
}
