package schemaadapter

import (
	"fmt"
	schemadomain "github.com/acme/streamforge-cdc/internal/schema/domain"
	"regexp"
	"strings"
)

var columnRE = regexp.MustCompile(`(?i)^([a-zA-Z0-9_]+)\s+([a-zA-Z0-9()]+)(?:\s+(NOT NULL))?`)

func ParseDDL(schema, table, ddl string, version int64) (schemadomain.Table, error) {
	ddl = strings.TrimSpace(ddl)
	start := strings.Index(ddl, "(")
	end := strings.LastIndex(ddl, ")")
	if start < 0 || end < start {
		return schemadomain.Table{}, fmt.Errorf("unsupported DDL")
	}
	var cols []schemadomain.Column
	for i, p := range strings.Split(ddl[start+1:end], ",") {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(strings.ToUpper(p), "PRIMARY KEY") {
			continue
		}
		m := columnRE.FindStringSubmatch(p)
		if len(m) == 0 {
			continue
		}
		cols = append(cols, schemadomain.Column{Name: m[1], Type: strings.ToUpper(m[2]), Nullable: m[3] == "", Position: i})
	}
	if len(cols) == 0 {
		return schemadomain.Table{}, fmt.Errorf("no columns found")
	}
	return schemadomain.Table{Schema: schema, Name: table, Version: version, Columns: cols}, nil
}
