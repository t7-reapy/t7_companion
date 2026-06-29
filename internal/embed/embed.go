// Package embed wraps go-sentex (pure-Go all-MiniLM-L6-v2) so the corpus build
// and the live query share ONE embedding implementation — exact-parity vectors.
package embed

import (
	"fmt"

	sentex "github.com/edgetools/go-sentex"
)

// Embedder produces 384-dim, L2-normalized sentence embeddings.
type Embedder struct {
	model *sentex.Model
}

// New loads the model (downloads ~87 MB on first use, cached thereafter).
func New() (*Embedder, error) {
	m, err := sentex.LoadModel()
	if err != nil {
		return nil, fmt.Errorf("load embedding model: %w", err)
	}
	return &Embedder{model: m}, nil
}

// Dim is the embedding dimension (384 for all-MiniLM-L6-v2).
func (e *Embedder) Dim() int { return e.model.Dimensions() }

// Embed returns the L2-normalized vector for text.
func (e *Embedder) Embed(text string) ([]float32, error) {
	return e.model.Embed(text)
}
