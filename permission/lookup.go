package permission

import "time"

const (
	MISS = iota
	PRIVATE
	PUBLIC
)

type LookupResponse struct {
	Public        bool
	Miss          bool
	LastCheckTime time.Time
}

type Lookup interface {
	Name() string
	CanView(string) (*LookupResponse, error)
}
