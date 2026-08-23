package logging

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Logger struct{ base *log.Logger }

func New() *Logger { return &Logger{base: log.New(os.Stdout, "", 0)} }
func (l *Logger) Log(level, msg string, fields map[string]any) {
	m := map[string]any{"time": time.Now().UTC().Format(time.RFC3339Nano), "level": level, "msg": msg}
	for k, v := range fields {
		m[k] = v
	}
	b, _ := json.Marshal(m)
	l.base.Println(string(b))
}
func (l *Logger) Info(msg string, f map[string]any)  { l.Log("info", msg, f) }
func (l *Logger) Error(msg string, f map[string]any) { l.Log("error", msg, f) }
