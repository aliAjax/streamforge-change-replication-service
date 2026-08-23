package bufferadapter

import (
	"bufio"
	"encoding/json"
	"fmt"
	buffer "github.com/acme/streamforge-cdc/internal/buffer/domain"
	"os"
	"sync"
)

type FileLog struct {
	mu   sync.Mutex
	path string
	file *os.File
}

func Open(path string) (*FileLog, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open buffer log: %w", err)
	}
	return &FileLog{path: path, file: f}, nil
}
func (l *FileLog) Append(r buffer.Record) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	if _, err = l.file.Write(append(b, '\n')); err != nil {
		return fmt.Errorf("write buffer log: %w", err)
	}
	return l.file.Sync()
}
func (l *FileLog) ReadAll() ([]buffer.Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.Open(l.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []buffer.Record
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4096), 16<<20)
	for scanner.Scan() {
		var r buffer.Record
		if err := json.Unmarshal(scanner.Bytes(), &r); err != nil {
			return nil, fmt.Errorf("decode buffer log: %w", err)
		}
		out = append(out, r)
	}
	return out, scanner.Err()
}
func (l *FileLog) Close() error { l.mu.Lock(); defer l.mu.Unlock(); return l.file.Close() }
