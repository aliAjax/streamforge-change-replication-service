package mysqladapter

import (
	"fmt"
	"strings"
)

type GTIDSet struct{ Values map[string][][2]uint64 }

func ParseGTID(v string) (GTIDSet, error) {
	out := GTIDSet{Values: map[string][][2]uint64{}}
	for _, part := range strings.Split(v, ",") {
		bits := strings.Split(part, ":")
		if len(bits) < 2 {
			return out, fmt.Errorf("invalid GTID %s", part)
		}
		uuid := bits[0]
		for _, r := range bits[1:] {
			var a, b uint64
			if _, err := fmt.Sscanf(r, "%d-%d", &a, &b); err != nil {
				if _, err = fmt.Sscanf(r, "%d", &a); err != nil {
					return out, err
				}
				b = a
			}
			out.Values[uuid] = append(out.Values[uuid], [2]uint64{a, b})
		}
	}
	return out, nil
}
func (s GTIDSet) Contains(uuid string, n uint64) bool {
	for _, r := range s.Values[uuid] {
		if n >= r[0] && n <= r[1] {
			return true
		}
	}
	return false
}
