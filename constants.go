package main

const (
	MTU = 1600
)

var tasks = []string{
	"a", "aaaa",
	"ds", "dsmtu",
	"v6mtu", "v6ns",
}

var nameMap = map[string]string{
	"a":     "Test with IPv4 DNS record",
	"aaaa":  "Test with IPv6 DNS record",
	"ds":    "Test with Dual Stack DNS record",
	"dsmtu": "Test for Dual Stack DNS and large packet",
	"v6mtu": "Test IPv6 large packet",
	"v6ns":  "Test if your ISP's DNS server uses IPv6",
	"asn4":  "Find IPv4 Service Provider",
	"asn6":  "Find IPv6 Service Provider",
}

var expectedTypes = map[string]string{
	"a":     "ipv4",
	"aaaa":  "ipv6",
	"ds":    "ipv6",
	"dsmtu": "ipv6",
	"v6mtu": "ipv6",
	"v6ns":  "ipv6",
	"asn4":  "ipv4",
	"asn6":  "ipv6",
}
