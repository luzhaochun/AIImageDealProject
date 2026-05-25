package idgen

import (
	"fmt"
	"sync"
	"time"
)

var seqMu sync.Mutex
var seq int64

func GlobalNo(tenantSlug string) string {
	seqMu.Lock()
	seq++
	n := seq
	seqMu.Unlock()
	return fmt.Sprintf("%s-%s-%06d", tenantSlug, time.Now().Format("20060102"), n%1000000)
}

func NewUUID() string {
	// use google uuid in callers; placeholder if needed
	return ""
}
