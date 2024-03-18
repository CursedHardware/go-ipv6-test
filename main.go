package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/CursedHardware/go-ipv6-test/ipv6test"
	"github.com/fatih/color"
	"net/http"
	"strings"
	"sync"
)

var host string
var listAll bool
var testAll bool

func init() {
	flag.StringVar(&host, "host", knownHosts[0], "Host")
	flag.BoolVar(&testAll, "all", false, "Test All Hosts")
	flag.BoolVar(&listAll, "hosts", false, "Available Hosts")
	flag.Parse()
}

func main() {
	tester := &ipv6test.Tester{
		Client: http.DefaultClient,
		MTU:    1600,
	}
	switch {
	case listAll:
		fmt.Println(strings.Join(knownHosts, "\n"))
	case testAll:
		var wg sync.WaitGroup
		var mux sync.Mutex
		for _, knownHost := range knownHosts {
			wg.Add(1)
			go func(host string) {
				report := tester.Run(ipv6test.RecordIPv6, host)
				mux.Lock()
				fmt.Printf("Test for %q\n", host)
				emitReport(report)
				mux.Unlock()
				wg.Done()
			}(knownHost)
		}
		wg.Wait()
	default:
		for _, taskType := range tasks {
			emitReport(tester.Run(taskType, host))
		}
		emitReport(tester.Run(ipv6test.RecordASN4, "ipv4.lookup.test-ipv6.com"))
		emitReport(tester.Run(ipv6test.RecordASN6, "ipv6.lookup.test-ipv6.com"))
	}
}

func emitReport(r *ipv6test.Report) {
	var ok = color.GreenString("ok")
	var bad = color.BlueString("bad")
	fmt.Println(r.Task.Name())
	if r.Failed {
		fmt.Printf("%s (%s)\n", bad, r.Elapsed)
		return
	}
	status := bad
	if r.Task.Match(r.IPProto) {
		status = ok
	}
	fmt.Printf("%s (%s) using %s", status, r.Elapsed, r.IPProto)
	if r.ASN != "" {
		fmt.Printf(" with ASN%s", r.ASN)
	}
	fmt.Println()
	fmt.Println()
}
