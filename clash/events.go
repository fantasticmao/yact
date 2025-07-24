package clash

type EventType int

const (
	EventTraffic EventType = 1 << iota
	EventTracing
)

type Traffic struct {
	Up   int64 `json:"up"`
	Down int64 `json:"down"`
}

type Basic struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Metadata struct {
	NetWork string `json:"network"`
	Type    string `json:"type"`
	SrcIP   string `json:"sourceIP"`
	DstIP   string `json:"destinationIP"`
	SrcPort string `json:"sourcePort"`
	DstPort string `json:"destinationPort"`
	Host    string `json:"host"`
	DNSMode string `json:"dnsMode"`
}

type EventRuleMatch struct {
	Basic
	Duration int64    `json:"duration"`
	Proxy    string   `json:"proxy"`
	Rule     string   `json:"rule"`
	Payload  string   `json:"payload"`
	Error    string   `json:"error"`
	Metadata Metadata `json:"metadata"`
}

type EventProxyDial struct {
	Basic
	Duration int64    `json:"duration"`
	Proxy    string   `json:"proxy"`
	Chain    []string `json:"chain"`
	Error    string   `json:"error"`
	Address  string   `json:"address"`
	Host     string   `json:"host"`
}

type EventDNSRequest struct {
	Basic
	Duration int64    `json:"duration"`
	Name     string   `json:"name"`
	Answer   []string `json:"answer"`
	Error    string   `json:"error"`
	QType    string   `json:"qType"`
	DNSType  string   `json:"dnsType"`
}
