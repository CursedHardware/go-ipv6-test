package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	. "github.com/CursedHardware/go-ipv6-test/ipv6test"
	"github.com/fatih/color"
)

var host string
var mirrors bool

func init() {
	flag.StringVar(&host, "host", "main.test-ipv6.com", "Host")
	flag.BoolVar(&mirrors, "mirrors", false, "List of all mirror sites")
	flag.Parse()
}

func main() {
	tester := &Tester{
		Client: http.DefaultClient,
		MTU:    1600,
	}
	switch {
	case mirrors:
		sites, err := tester.GetKnownSites()
		if err != nil {
			log.Panicln(err)
			return
		}
		for _, site := range sites {
			if site.Hidden || !site.Mirror {
				continue
			}
			fmt.Println(site.FQDN)
		}
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
