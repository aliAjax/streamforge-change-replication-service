package schemadomain

import (
	"fmt"
	"strings"
	"time"
)

type Strategy string

const (
	Strict     Strategy = "strict"
	Compatible Strategy = "compatible"
	Raw        Strategy = "raw"
)

type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Position int    `json:"position"`
}
type Table struct {
	Schema     string    `json:"schema"`
	Name       string    `json:"name"`
	Version    int64     `json:"version"`
	PrimaryKey []string  `json:"primary_key"`
	Columns    []Column  `json:"columns"`
	Hash       string    `json:"hash"`
	CapturedAt time.Time `json:"captured_at"`
}

func (t Table) Column(name string) (Column, bool) {
	for _, c := range t.Columns {
		if c.Name == name {
			return c, true
		}
	}
	return Column{}, false
}
func Compare(old, new Table) []string {
	var changes []string
	oldMap := map[string]Column{}
	for _, c := range old.Columns {
		oldMap[c.Name] = c
	}
	newMap := map[string]Column{}
	for _, c := range new.Columns {
		newMap[c.Name] = c
		if p, ok := oldMap[c.Name]; !ok {
			changes = append(changes, "added:"+c.Name)
		} else if p.Type != c.Type || p.Nullable != c.Nullable {
			changes = append(changes, fmt.Sprintf("changed:%s:%s->%s", c.Name, p.Type, c.Type))
		}
	}
	for _, c := range old.Columns {
		if _, ok := newMap[c.Name]; !ok {
			changes = append(changes, "removed:"+c.Name)
		}
	}
	if strings.Join(old.PrimaryKey, ",") != strings.Join(new.PrimaryKey, ",") {
		changes = append(changes, "primary_key_changed")
	}
	return changes
}
