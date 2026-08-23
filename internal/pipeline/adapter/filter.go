package pipelineadapter

import (
	transaction "github.com/acme/streamforge-cdc/internal/transaction/domain"
	"strings"
)

type Filter struct {
	Tables  map[string]bool
	Columns map[string]map[string]bool
}

func (f Filter) Accept(e transaction.Event) bool {
	if len(f.Tables) > 0 && !f.Tables[e.Schema+"."+e.Table] && !f.Tables[e.Table] {
		return false
	}
	return true
}
func (f Filter) Project(e transaction.Event) transaction.Event {
	if cols, ok := f.Columns[e.Schema+"."+e.Table]; ok {
		e.Before = project(e.Before, cols)
		e.After = project(e.After, cols)
	}
	return e
}
func project(in map[string]any, cols map[string]bool) map[string]any {
	if in == nil {
		return nil
	}
	out := map[string]any{}
	for k, v := range in {
		if cols[k] || strings.HasPrefix(k, "__") {
			out[k] = v
		}
	}
	return out
}
