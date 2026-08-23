package sourcedomain

import "time"

type Kind string

const (
	PostgreSQL Kind = "postgres"
	MySQL      Kind = "mysql"
)

type Status string

const (
	Draft      Status = "draft"
	Validating Status = "validating"
	Active     Status = "active"
	Paused     Status = "paused"
	Failed     Status = "failed"
	Retired    Status = "retired"
)

type Filter struct {
	Schemas         []string            `json:"schemas"`
	Tables          []string            `json:"tables"`
	ExcludedColumns map[string][]string `json:"excluded_columns"`
}
type Source struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Kind          Kind      `json:"kind"`
	Endpoint      string    `json:"endpoint"`
	Database      string    `json:"database"`
	Filter        Filter    `json:"filter"`
	Status        Status    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CredentialRef string    `json:"credential_ref"`
}

func (s Source) Allows(schema, table string) bool {
	for _, x := range s.Filter.Schemas {
		if x == schema {
			goto tableCheck
		}
	}
	if len(s.Filter.Schemas) > 0 {
		return false
	}
tableCheck:
	for _, x := range s.Filter.Tables {
		if x == table {
			return true
		}
	}
	return len(s.Filter.Tables) == 0
}
