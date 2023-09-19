package tracingstruct

type LogType struct {
	Type string `json:"type"`
}

type Traffic struct {
	Up         int   `json:"up"`
	Down       int   `json:"down"`
	CreateTime int64 `json:"createTime"`
}

type RuleMatch struct {
	Dnstime  int    `json:"dnstime"`
	Duration int    `json:"duration"`
	ID       string `json:"id"`
	Metadata struct {
		Network         string `json:"network"`
		Type            string `json:"type"`
		SourceIP        string `json:"sourceIP"`
		DestinationIP   string `json:"destinationIP"`
		SourcePort      string `json:"sourcePort"`
		DestinationPort string `json:"destinationPort"`
		Host            string `json:"host"`
		DNSMode         string `json:"dnsMode"`
		ProcessPath     string `json:"processPath"`
		SpecialProxy    string `json:"specialProxy"`
	} `json:"metadata"`
	Payload    string `json:"payload"`
	Proxy      string `json:"proxy"`
	Rule       string `json:"rule"`
	Type       string `json:"type"`
	CreateTime int64  `json:"createTime"`
}

type ProxyDial struct {
	Address    string   `json:"address"`
	Chain      []string `json:"chain"`
	Duration   int      `json:"duration"`
	Host       string   `json:"host"`
	ID         string   `json:"id"`
	Proxy      string   `json:"proxy"`
	Type       string   `json:"type"`
	CreateTime int64    `json:"createTime"`
}

type DNSRequest struct {
	Answer     []string `json:"answer"`
	DNSType    string   `json:"dnsType"`
	Duration   int      `json:"duration"`
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	QType      string   `json:"qType"`
	Source     string   `json:"source"`
	Type       string   `json:"type"`
	CreateTime int64    `json:"createTime"`
}
