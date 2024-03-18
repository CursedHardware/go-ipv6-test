package ipv6test

type Task uint8

const (
	RecordIPv4 Task = iota + 1
	RecordASN4
	RecordIPv6
	RecordDualStack
	RecordDualStackMTU
	RecordIPv6MTU
	RecordIPv6NS
	RecordASN6
)

func (t Task) Name() string {
	switch t {
	case RecordIPv4:
		return "Test with IPv4 DNS record"
	case RecordIPv6:
		return "Test with IPv6 DNS record"
	case RecordDualStack:
		return "Test with Dual Stack DNS record"
	case RecordDualStackMTU:
		return "Test for Dual Stack DNS and large packet"
	case RecordIPv6MTU:
		return "Test IPv6 large packet"
	case RecordIPv6NS:
		return "Test if your ISP's DNS server uses IPv6"
	case RecordASN4:
		return "Find IPv4 Service Provider"
	case RecordASN6:
		return "Find IPv6 Service Provider"
	}
	return ""
}

func (t Task) Prefix() string {
	switch t {
	case RecordIPv4, RecordASN4:
		return "ipv4."
	case RecordIPv6, RecordASN6:
		return "ipv6."
	case RecordDualStack, RecordDualStackMTU:
		return "ds."
	case RecordIPv6MTU:
		return "mtu1280."
	case RecordIPv6NS:
		return "ds.v6ns."
	}
	return ""
}

func (t Task) Match(proto string) bool {
	switch proto {
	case "ipv4":
		return t >= RecordIPv4 && t <= RecordASN4
	case "ipv6":
		return t >= RecordIPv6 && t <= RecordASN6
	}
	return false
}
