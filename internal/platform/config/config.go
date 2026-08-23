package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr        string
	MetricsAddr     string
	BufferBytes     int
	ShutdownTimeout time.Duration
	APIKey          string
}

func Load() Config {
	c := Config{HTTPAddr: ":8092", MetricsAddr: ":8093", BufferBytes: 8 << 20, ShutdownTimeout: 10 * time.Second, APIKey: "dev-key"}
	if v := os.Getenv("STREAMFORGE_HTTP_ADDR"); v != "" {
		c.HTTPAddr = v
	}
	if v := os.Getenv("STREAMFORGE_API_KEY"); v != "" {
		c.APIKey = v
	}
	if v := os.Getenv("STREAMFORGE_BUFFER_BYTES"); v != "" {
		if n, e := strconv.Atoi(v); e == nil && n > 0 {
			c.BufferBytes = n
		}
	}
	return c
}
