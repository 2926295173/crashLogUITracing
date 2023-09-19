package tracingstruct

type DNSCounting struct {
	Name     string  `json:"name"`
	Duration float64 `json:"duration"`
	Count    int     `json:"count"`
}

type DNSQTypeCounting struct {
	QType    string  `json:"qType"`
	Duration float64 `json:"duration"`
	Count    int     `json:"count"`
}

type DNSTypeCounting struct {
	DNSType  string  `json:"dnsType"`
	Duration float64 `json:"duration"`
	Count    int     `json:"count"`
}

type ApiDnsRequest struct {
	Counting        []DNSCounting      `json:"counting"`
	QTypeCounting   []DNSQTypeCounting `json:"qTypeCounting"`
	DNSTypeCounting []DNSTypeCounting  `json:"dnsTypeCounting"`
}

type RuleMatchCounting struct {
	Duration float64 `json:"duration"`
	Payload  string  `json:"payload"`
	Proxy    string  `json:"proxy"`
	Count    int     `json:"count"`
}

type RuleMatchPortCounting struct {
	Port  int `json:"port"`
	Count int `json:"count"`
}

type RuleMatchProcessCounting struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

type RuleMatchClientCounting struct {
	IP    string `json:"ip"`
	Count int    `json:"count"`
}

type ApiRuleMatchRequest struct {
	Counting        []RuleMatchCounting        `json:"counting"`
	PortCounting    []RuleMatchPortCounting    `json:"portCounting"`
	ProcessCounting []RuleMatchProcessCounting `json:"processCounting"`
	ClientCounting  []RuleMatchClientCounting  `json:"clientCounting"`
}

type ProxyDialProxyCounting struct {
	Proxy    string  `json:"proxy"`
	Duration float64 `json:"duration"`
	Count    int     `json:"count"`
}

type ProxyDialHostCounting struct {
	Host     string  `json:"host"`
	Duration float64 `json:"duration"`
	Count    int     `json:"count"`
}

type ApiProxyDialRequest struct {
	ProxyCounting []ProxyDialProxyCounting `json:"proxyCounting"`
	HostCounting  []ProxyDialHostCounting  `json:"hostCounting"`
}

type ApiTrafficHistory struct {
	Up         float64 `json:"up"`
	Down       float64 `json:"down"`
	CreateTime int     `json:"createTime"`
}

type ApiTrafficRequest struct {
	Up      float64             `json:"up"`
	Down    float64             `json:"down"`
	History []ApiTrafficHistory `json:"history"`
}
