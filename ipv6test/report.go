package ipv6test

import (
	"time"
)

type Report struct {
	Task    Task
	Host    string
	IPProto string
	ASN     string
	Failed  bool
	Elapsed time.Duration
}
