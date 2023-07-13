package core

import (
	"time"
)

type Token struct {
	Token   string
	Expired time.Time
}

func (t *Token) Validate() bool {
	return time.Now().Before(t.Expired.Add(5 * time.Minute))
}
