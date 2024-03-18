package ipv6test

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Tester struct {
	Client *http.Client
	MTU    int
}

func (t *Tester) Run(task Task, host string) (report *Report) {
	u := &url.URL{Scheme: "https", Host: task.Prefix() + host, Path: "/ip/"}
	request, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0")
	if task == RecordIPv6MTU || task == RecordDualStackMTU {
		request.URL.RawQuery = "?fill=" + strings.Repeat("x", t.MTU)
	}
	start := time.Now()
	report = &Report{Task: task, Host: host}
	response, err := t.Client.Do(request)
	report.Elapsed = time.Since(start)
	if err != nil {
		report.Failed = true
		return
	}
	result := new(struct {
		Type string `json:"type"`
		ASN  string `json:"asn"`
	})
	if err = unmarshalJSONP(response.Body, result); err != nil {
		report.Failed = true
	}
	report.IPProto = result.Type
	report.ASN = result.ASN
	return
}
