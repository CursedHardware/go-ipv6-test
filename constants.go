package main

import "github.com/CursedHardware/go-ipv6-test/ipv6test"

var knownHosts = []string{
	"jp2.test-ipv6.com",
	"ams2.test-ipv6.com",
	"fra.test-ipv6.com",
	"sixte.st",
	"test-ipv6.arauc.br",
	"test-ipv6.belwue.net",
	"test-ipv6.carnet.hr",
	"test-ipv6.cl",
	"test-ipv6.coloradointerlink.net",
	"test-ipv6.com.au",
	"test-ipv6.cs.umd.edu",
	"test-ipv6.csclub.uwaterloo.ca",
	"test-ipv6.cz",
	"test-ipv6.epic.network",
	"test-ipv6.fratec.net",
	"test-ipv6.freerangecloud.com",
	"test-ipv6.go6.si",
	"test-ipv6.hkg.vr.org",
	"test-ipv6.hu",
	"test-ipv6.is",
	"test-ipv6.iu13.net",
	"test-ipv6.ke.liquidtelecom.net",
	"test-ipv6.noroutetohost.net",
	"test-ipv6.roedu.net",
	"test-ipv6.se",
	"test-ipv6.sin.vr.org",
	"test-ipv6.ttk.ru",
	"test-ipv6.vzxy.net",
	"testipv6.cn",
	"testipv6.de",
}

var tasks = []ipv6test.Task{
	ipv6test.RecordIPv4,
	ipv6test.RecordIPv6,
	ipv6test.RecordDualStack,
	ipv6test.RecordDualStackMTU,
	ipv6test.RecordIPv6MTU,
	ipv6test.RecordIPv6NS,
}
