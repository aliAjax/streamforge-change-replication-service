package checkpointadapter

import (
	"encoding/json"
	"fmt"
	checkpoint "github.com/acme/streamforge-cdc/internal/checkpoint/domain"
	"os"
	"path/filepath"
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
	if err != nil { return err }
	dir := filepath.Dir(f.path)
	tmp, err := os.CreateTemp(dir, ".checkpoint-")
	if err != nil { return err }
	if _, err = tmp.Write(b); err != nil { _ = tmp.Close(); return err }
	if err = tmp.Close(); err != nil { return err }
	if err = os.Rename(tmp.Name(), f.path); err != nil { return fmt.Errorf("checkpoint rename: %w", err) }
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
