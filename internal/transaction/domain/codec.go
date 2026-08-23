package transactiondomain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

func EncodeValue(v any) any {
	switch x := v.(type) {
	case []byte:
		return map[string]string{"$binary": base64.StdEncoding.EncodeToString(x)}
	case time.Time:
		return map[string]string{"$time": x.UTC().Format(time.RFC3339Nano)}
	case *big.Int:
		return map[string]string{"$bigint": x.String()}
	case json.Number:
		return map[string]string{"$decimal": x.String()}
	default:
		return v
	}
}

func DecodeValue(v any) (any, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return v, nil
	}
	if b, ok := m["$binary"].(string); ok {
		out, err := base64.StdEncoding.DecodeString(b)
		if err != nil {
			return nil, fmt.Errorf("decode binary: %w", err)
		}
		return out, nil
	}
	if t, ok := m["$time"].(string); ok {
		out, err := time.Parse(time.RFC3339Nano, t)
		if err != nil {
			return nil, fmt.Errorf("decode time: %w", err)
		}
		return out, nil
	}
	if n, ok := m["$bigint"].(string); ok {
		out, ok := new(big.Int).SetString(n, 10)
		if !ok {
			return nil, fmt.Errorf("decode bigint %s", n)
		}
		return out, nil
	}
	if d, ok := m["$decimal"].(string); ok {
		if _, err := strconv.ParseFloat(d, 64); err != nil {
			return nil, fmt.Errorf("decode decimal %s", d)
		}
		return json.Number(d), nil
	}
	return v, nil
}
