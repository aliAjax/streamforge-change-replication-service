package checkpointadapter

import (
	"encoding/json"
	"fmt"
	checkpoint "github.com/acme/streamforge-cdc/internal/checkpoint/domain"
	"os"
	"sync"
)

type FileStore struct {
	mu   sync.Mutex
	path string
}

func New(path string) *FileStore { return &FileStore{path: path} }
func (f *FileStore) Save(states []checkpoint.State) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	b, err := json.MarshalIndent(states, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(f.path, b, 0600); err != nil {
		return fmt.Errorf("checkpoint save: %w", err)
	}
	return nil
}
func (f *FileStore) Load() ([]checkpoint.State, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	b, err := os.ReadFile(f.path)
	if os.IsNotExist(err) {
		return []checkpoint.State{}, nil
	}
	if err != nil {
		return nil, err
	}
	var out []checkpoint.State
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("checkpoint decode: %w", err)
	}
	return out, nil
}
