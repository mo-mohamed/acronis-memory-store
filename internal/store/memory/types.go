package memory

import "time"

type Value struct {
	Val    string
	TTL    time.Time
	IsList bool
	List   []string
}
