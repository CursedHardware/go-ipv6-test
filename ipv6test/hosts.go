package ipv6test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

type TestSite struct {
	FQDN        string `json:"site"`
	Mirror      bool   `json:"mirror"`
	Hidden      bool   `json:"hide"`
	Location    string `json:"loc"`
	Provider    string `json:"provider"`
	Monitor     string `json:"monitor"`
	Contact     string `json:"contact"`
	Reason      string `json:"reason"`
	Transparent bool   `json:"transparent"`
}

func (t *Tester) GetKnownSites() (knownSites map[string]TestSite, err error) {
	response, err := t.Client.Get("https://www.test-ipv6.com/index.js.gz.en_US")
	if err != nil {
		return
	}
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	const prefix = "GIGO.sites_parsed={"
	const suffix = "};"
	var startIndex, endIndex int
	startIndex = bytes.Index(payload, []byte(prefix))
	if startIndex == -1 {
		return nil, errors.New("ipv6test: unable to find known sites")
	}
	startIndex += len(prefix) - 1
	endIndex = bytes.Index(payload[startIndex:], []byte(suffix))
	if endIndex == -1 {
		return nil, errors.New("ipv6test: unable to find known sites")
	}
	endIndex += startIndex + len(suffix) - 1
	payload = payload[startIndex:endIndex]
	err = json.Unmarshal(payload, &knownSites)
	return
}
