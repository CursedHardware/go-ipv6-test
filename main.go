package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ok = color.GreenString("ok")
var bad = color.BlueString("bad")

//go:embed tests.json
var bundle []byte
var tests map[string]map[string]string

var host string
var listAll bool
var testAll bool

func init() {
	tests = make(map[string]map[string]string)
	if err := json.Unmarshal(bundle, &tests); err != nil {
		panic(err)
	}
	flag.StringVar(&host, "host", "test-ipv6.com", "Host")
	flag.BoolVar(&testAll, "all", false, "Test All Sites")
	flag.BoolVar(&listAll, "sites", false, "Available Sites")
	flag.Parse()
}

func main() {
	switch {
	case testAll:
		for _, testHost := range getAllHosts() {
			fmt.Printf("Test for %q\n", testHost)
			run(testHost, "aaaa")
		}
	case listAll:
		fmt.Println(strings.Join(getAllHosts(), "\n"))
		return
	default:
		for _, taskName := range tasks {
			run(host, taskName)
		}
	}
	runWithHost("ipv4.lookup.test-ipv6.com", "test-ipv6.com", "asn4")
	runWithHost("ipv6.lookup.test-ipv6.com", "test-ipv6.com", "asn4")
}

func run(testHost, taskName string) {
	runWithHost(tests[testHost][taskName], testHost, taskName)
}

func runWithHost(host, testHost, taskName string) {
	start := time.Now()
	report, err := do(taskName, host, testHost, http.DefaultClient)
	elapsed := time.Since(start)
	fmt.Println(nameMap[taskName])
	if err != nil {
		fmt.Printf("%s (%s)\n", bad, elapsed)
	} else if report.ASN != "" {
		status := bad
		if expectedTypes[taskName] == report.Type {
			status = ok
		}
		fmt.Printf("%s (%s) using %s with ASN%s\n", status, elapsed, report.Type, report.ASN)
	} else {
		status := bad
		if expectedTypes[taskName] == report.Type {
			status = ok
		}
		fmt.Printf("%s (%s) using %s\n", status, elapsed, report.Type)
	}
	fmt.Println()
}

func getAllHosts() []string {
	var names []string
	for name := range tests {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func do(name, host, testHost string, client *http.Client) (report *Report, err error) {
	request, _ := http.NewRequest(http.MethodGet, buildURL(name, host, testHost).String(), nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0")
	response, err := client.Do(request)
	if err != nil {
		return
	}
	report = new(Report)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	start := bytes.IndexByte(body, '{')
	end := bytes.IndexByte(body, '}') + 1
	err = json.Unmarshal(body[start:end], &report)
	return
}

func buildURL(name, host, testHost string) *url.URL {
	query := make(url.Values)
	query.Set("testname", fmt.Sprint("test_", name))
	query.Set("testdomain", host)
	switch {
	case strings.HasSuffix(name, "mtu"):
		query.Set("size", strconv.Itoa(MTU))
		query.Set("fill", strings.Repeat("x", MTU))
	case strings.HasPrefix(name, "asn"):
		query.Set("asn", "1")
	}
	return &url.URL{
		Scheme:   "https",
		Host:     host,
		Path:     "/ip/",
		RawQuery: query.Encode(),
	}
}
