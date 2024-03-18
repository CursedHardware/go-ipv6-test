package main

import (
	_ "embed"
	"flag"
	"fmt"
	. "github.com/CursedHardware/go-ipv6-test/ipv6test"
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
	tester := &Tester{
		Client: http.DefaultClient,
		MTU:    1600,
	}
	switch {
	case listAll:
		fmt.Println(strings.Join(knownHosts, "\n"))
	case testAll:
		batchTasks(tester, []Task{RecordIPv6}, knownHosts)
	default:
		fullTask(tester, host)
	}
}

func invoke(tester *Tester, requests map[string][]Task) <-chan *Report {
	var wg sync.WaitGroup
	reports := make(chan *Report)
	for taskHost, tasks := range requests {
		wg.Add(len(tasks))
		for _, task := range tasks {
			go func(task Task, taskHost string) {
				if task == RecordASN4 || task == RecordASN6 {
					taskHost = "lookup.test-ipv6.com"
				}
				reports <- tester.Run(task, taskHost)
				wg.Done()
			}(task, taskHost)
		}
	}
	go func() {
		wg.Wait()
		close(reports)
	}()
	return reports
}

func fullTask(tester *Tester, host string) {
	requests := make(map[string][]Task)
	requests[host] = []Task{
		RecordIPv4, RecordIPv6, RecordDualStack, RecordDualStackMTU,
		RecordIPv6MTU, RecordIPv6NS, RecordASN4, RecordASN6,
	}
	for report := range invoke(tester, requests) {
		emitReport(report)
	}
}

func batchTasks(tester *Tester, tasks []Task, hosts []string) {
	requests := make(map[string][]Task)
	for _, testHost := range hosts {
		requests[testHost] = tasks
	}
	for report := range invoke(tester, requests) {
		fmt.Printf("Test for %q\n", report.Host)
		emitReport(report)
	}
}

func emitReport(r *Report) {
	var ok = color.GreenString("ok")
	var bad = color.BlueString("bad")
	fmt.Println(r.Task.Name())
	if r.Failed {
		fmt.Printf("%s (%s)\n", bad, r.Elapsed)
		fmt.Println()
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
