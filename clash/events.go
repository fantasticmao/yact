package clash

type EventMode int

const (
	EventModeTraffic EventMode = 1 << iota
	EventModeTracing
)

type EventTraffic struct {
	Up   int64 `json:"up"`
	Down int64 `json:"down"`
}

const (
	EventTypeRuleMatch  string = "RuleMatch"
	EventTypeProxyDial  string = "ProxyDial"
	EventTypeDNSRequest string = "DNSRequest"
)

type EventBasic struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Duration int64  `json:"duration"`
	Error    string `json:"error"`
}

type EventMetadata struct {
	NetWork string `json:"network"`
	Type    string `json:"type"`
	SrcIP   string `json:"sourceIP"`
	DstIP   string `json:"destinationIP"`
	SrcPort string `json:"sourcePort"`
	DstPort string `json:"destinationPort"`
	Host    string `json:"host"`
	DnsMode string `json:"dnsMode"`
}

type EventRuleMatch struct {
	EventBasic
	Proxy    string        `json:"proxy"`
	Rule     string        `json:"rule"`
	Payload  string        `json:"payload"`
	Metadata EventMetadata `json:"metadata"`
}

type EventProxyDial struct {
	EventBasic
	Proxy   string   `json:"proxy"`
	Chain   []string `json:"chain"`
	Address string   `json:"address"`
	Host    string   `json:"host"`
}

type EventDnsRequest struct {
	EventBasic
	DndType string   `json:"dnsType"`
	Name    string   `json:"name"`
	QType   string   `json:"qType"`
	Answer  []string `json:"answer"`
	Source  string   `json:"source"`
}
