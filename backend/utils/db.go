package utils

import (
	"clash-tracing/tracingstruct"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	mutex            sync.Mutex
	bufferSize       = 64
	interval         = 64
	latestTraffic    tracingstruct.Traffic
	dnsRequestBuffer = []tracingstruct.DNSRequest{}
	ruleMatchBuffer  = []tracingstruct.RuleMatch{}
	proxyDialBuffer  = []tracingstruct.ProxyDial{}
	trafficBuffer    = []tracingstruct.Traffic{}
)

func InitDBTables(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS "traffic" (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"up" integer,
		"down" integer,
		"createTime" integer
	);
	CREATE TABLE IF NOT EXISTS "proxy_dial" (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"address" TEXT,
		"duration" integer,
		"host" TEXT,
		"proxy" TEXT,
		"createTime" integer
	);
	CREATE TABLE IF NOT EXISTS "dns_request" (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"answer" TEXT,
		"dnsType" TEXT,
		"duration" integer,
		"name" TEXT,
		"qType" TEXT,
		"createTime" integer
	);
	CREATE TABLE IF NOT EXISTS "rule_match" (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"duration" integer,
		"rule" TEXT,
		"payload" TEXT,
		"proxy" TEXT,
		"m_network" TEXT,
		"m_type" TEXT,
		"m_sourceIP" TEXT,
		"m_destinationIP" TEXT,
		"m_sourcePort" TEXT,
		"m_destinationPort" TEXT,
		"m_host" TEXT,
		"m_dnsMode" TEXT,
		"m_processPath" TEXT,
		"createTime" integer
	);
	`
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB() *sql.DB {
	v := url.Values{}
	v.Set("cache", "shared")
	v.Set("mode", "rwc")
	u := url.URL{
		Scheme:   "file",
		Opaque:   url.PathEscape("clash-tracing.db"),
		RawQuery: v.Encode(),
	}
	db, err := sql.Open("sqlite3", u.String())
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	fmt.Println("Connected to Local Database")

	return db
}

func HandleDNSRequestLog(msg []byte, db *sql.DB) {
	dns := tracingstruct.DNSRequest{
		CreateTime: time.Now().Unix(),
	}
	err := json.Unmarshal(msg, &dns)
	if err != nil {
		return
	}

	mutex.Lock()

	dnsRequestBuffer = append(dnsRequestBuffer, dns)

	if len(dnsRequestBuffer) >= bufferSize {
		syncDNSRequest(db)
	}

	mutex.Unlock()
}

func HandleRuleMatchLog(msg []byte, db *sql.DB) {
	rule := tracingstruct.RuleMatch{
		CreateTime: time.Now().Unix(),
	}
	err := json.Unmarshal(msg, &rule)
	if err != nil {
		return
	}

	mutex.Lock()

	ruleMatchBuffer = append(ruleMatchBuffer, rule)

	if len(ruleMatchBuffer) >= bufferSize {
		syncRuleMatch(db)
	}

	mutex.Unlock()
}

func HandleProxyDialLog(msg []byte, db *sql.DB) {
	proxy := tracingstruct.ProxyDial{
		CreateTime: time.Now().Unix(),
	}
	err := json.Unmarshal(msg, &proxy)
	if err != nil {
		return
	}

	mutex.Lock()

	proxyDialBuffer = append(proxyDialBuffer, proxy)

	if len(proxyDialBuffer) >= bufferSize {
		syncProxyDial(db)
	}

	mutex.Unlock()
}

func HandleTraffic(msg []byte, db *sql.DB) {
	traffic := tracingstruct.Traffic{
		CreateTime: time.Now().Unix(),
	}
	err := json.Unmarshal(msg, &traffic)
	if err != nil {
		return
	}

	mutex.Lock()

	if traffic.Up == 0 && traffic.Down == 0 {
		if latestTraffic.Up != 0 || latestTraffic.Down != 0 {
			trafficBuffer = append(trafficBuffer, traffic)
		}
	} else {
		trafficBuffer = append(trafficBuffer, traffic)
	}

	latestTraffic = traffic

	if len(trafficBuffer) >= interval {
		syncTraffic(db)
	}

	mutex.Unlock()
}

func startInsert(db *sql.DB, logMsg string, query string, callback func(*sql.Stmt, *sql.Tx)) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("[%s]\t Failed to start transaction:%v\n", logMsg, err)
		return
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("[%s]\t Failed to prepare statement:%v\n", logMsg, err)
		tx.Rollback()
		return
	}

	callback(stmt, tx)

	err = tx.Commit()

	if err != nil {
		log.Printf("[%s]\t Failed to commit transaction:%v\n", logMsg, err)
		return
	}
}

func QueryDnsRequests(db *sql.DB) tracingstruct.ApiDnsRequest {
	mutex.Lock()

	dnsRequest := tracingstruct.ApiDnsRequest{
		Counting:        []tracingstruct.DNSCounting{},
		QTypeCounting:   []tracingstruct.DNSQTypeCounting{},
		DNSTypeCounting: []tracingstruct.DNSTypeCounting{},
	}

	// 域名查询次数统计
	rows, _ := db.Query(`SELECT name, AVG(duration) as duration, COUNT(1) as count FROM "dns_request" GROUP BY name ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		r := tracingstruct.DNSCounting{}
		rows.Scan(&r.Name, &r.Duration, &r.Count)
		dnsRequest.Counting = append(dnsRequest.Counting, r)
	}

	// DNS queryType 统计
	rows, _ = db.Query(`SELECT qType, AVG(duration) as duration, COUNT(1) as count FROM "dns_request" GROUP BY qType ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		r := tracingstruct.DNSQTypeCounting{}
		rows.Scan(&r.QType, &r.Duration, &r.Count)
		dnsRequest.QTypeCounting = append(dnsRequest.QTypeCounting, r)
	}

	// DNS dnsType统计
	rows, _ = db.Query(`SELECT dnsType, AVG(duration) as duration, COUNT(1) as count FROM "dns_request" GROUP BY dnsType ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		r := tracingstruct.DNSTypeCounting{}
		rows.Scan(&r.DNSType, &r.Duration, &r.Count)
		dnsRequest.DNSTypeCounting = append(dnsRequest.DNSTypeCounting, r)
	}

	mutex.Unlock()

	return dnsRequest
}

func QueryRuleMatchRequests(db *sql.DB) tracingstruct.ApiRuleMatchRequest {
	mutex.Lock()

	ruleMatch := tracingstruct.ApiRuleMatchRequest{
		Counting:        []tracingstruct.RuleMatchCounting{},
		PortCounting:    []tracingstruct.RuleMatchPortCounting{},
		ProcessCounting: []tracingstruct.RuleMatchProcessCounting{},
		ClientCounting:  []tracingstruct.RuleMatchClientCounting{},
	}

	// 查询每种ruleset规则的匹配次数和平均时间
	rows, _ := db.Query(`SELECT AVG(duration) as duration, payload, proxy, COUNT(1) as count FROM "rule_match" GROUP BY payload ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		r := tracingstruct.RuleMatchCounting{}
		rows.Scan(&r.Duration, &r.Payload, &r.Proxy, &r.Count)
		ruleMatch.Counting = append(ruleMatch.Counting, r)
	}

	// 查询端口来匹配流量的类型
	rows, _ = db.Query(`SELECT m_destinationPort, COUNT(1) as count from rule_match GROUP BY m_destinationPort ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		r := tracingstruct.RuleMatchPortCounting{}
		rows.Scan(&r.Port, &r.Count)
		ruleMatch.PortCounting = append(ruleMatch.PortCounting, r)
	}

	// 查询路径来统计各个APP的连接数
	rows, _ = db.Query(`SELECT m_processPath, COUNT(1) as count FROM "rule_match" GROUP BY m_processPath ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		r := tracingstruct.RuleMatchProcessCounting{}
		rows.Scan(&r.Path, &r.Count)
		ruleMatch.ProcessCounting = append(ruleMatch.ProcessCounting, r)
	}

	// 统计各个客户端IP的请求次数
	rows, _ = db.Query(`SELECT m_sourceIP, COUNT(1) as count FROM "rule_match" GROUP BY m_sourceIP ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		r := tracingstruct.RuleMatchClientCounting{}
		rows.Scan(&r.IP, &r.Count)
		ruleMatch.ClientCounting = append(ruleMatch.ClientCounting, r)
	}

	mutex.Unlock()

	return ruleMatch
}

func QueryProxyDialRequests(db *sql.DB) tracingstruct.ApiProxyDialRequest {
	mutex.Lock()

	proxyDial := tracingstruct.ApiProxyDialRequest{
		ProxyCounting: []tracingstruct.ProxyDialProxyCounting{},
		HostCounting:  []tracingstruct.ProxyDialHostCounting{},
	}

	// 查询代理节点被使用的次数和时间
	rows, _ := db.Query(`SELECT proxy, AVG(duration) as duration, COUNT(1) as count FROM proxy_dial GROUP BY proxy ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		r := tracingstruct.ProxyDialProxyCounting{}
		rows.Scan(&r.Proxy, &r.Duration, &r.Count)
		proxyDial.ProxyCounting = append(proxyDial.ProxyCounting, r)
	}

	// 查询域名被访问次数和时间
	rows, _ = db.Query(`SELECT address, AVG(duration) as duration, COUNT(1) as count FROM proxy_dial GROUP BY address ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		r := tracingstruct.ProxyDialHostCounting{}
		rows.Scan(&r.Host, &r.Duration, &r.Count)
		proxyDial.HostCounting = append(proxyDial.HostCounting, r)
	}

	mutex.Unlock()

	return proxyDial
}

func QueryTrafficRequests(db *sql.DB) tracingstruct.ApiTrafficRequest {
	mutex.Lock()

	traffic := tracingstruct.ApiTrafficRequest{
		Up:      0,
		Down:    0,
		History: []tracingstruct.ApiTrafficHistory{},
	}

	// 查询平均速度
	db.QueryRow(`SELECT AVG(up) as up, AVG(down) as down FROM "traffic"`).Scan(&traffic.Up, &traffic.Down)

	// 查询流量历史
	rows, _ := db.Query(`SELECT up, down, createTime FROM "traffic" ORDER BY createTime DESC LIMIT 10`)
	defer rows.Close()

	for rows.Next() {
		r := tracingstruct.ApiTrafficHistory{}
		rows.Scan(&r.Up, &r.Down, &r.CreateTime)
		traffic.History = append(traffic.History, r)
	}

	mutex.Unlock()

	return traffic
}

func QueryProcessDetail(db *sql.DB, path string, page int, pageSize int) tracingstruct.PageType[[]tracingstruct.ApiProcessDetail] {
	mutex.Lock()

	res := tracingstruct.PageType[[]tracingstruct.ApiProcessDetail]{
		Total:    0,
		Page:     page,
		PageSize: pageSize,
		Data:     []tracingstruct.ApiProcessDetail{},
	}

	// 根据进程查询请求
	rows, _ := db.Query(`SELECT m_sourceIP, m_sourcePort, m_destinationIP, m_destinationPort, m_host, m_dnsMode, createTime FROM "rule_match" WHERE m_processPath = ? LIMIT ?,?`, path, (page-1)*pageSize, pageSize)
	defer rows.Close()

	for rows.Next() {
		r := tracingstruct.ApiProcessDetail{}
		rows.Scan(&r.SourceIP, &r.SourcePort, &r.DestinationIP, &r.DestinationPort, &r.Host, &r.DnsMode, &r.CreateTime)
		res.Data = append(res.Data, r)
	}

	// 查总数
	db.QueryRow(`SELECT COUNT(1) FROM "rule_match" where m_processPath = ?`, path).Scan(&res.Total)

	mutex.Unlock()

	return res
}

func syncDNSRequest(db *sql.DB) {
	logMsg, query := "DNSRequest", "INSERT INTO dns_request(answer,dnsType,duration,name,qType, createTime) values(?,?,?,?,?,?)"
	fmt.Printf("[%s]\t Writing %d records to the database...\n", logMsg, len(dnsRequestBuffer))
	startInsert(db, logMsg, query, func(stmt *sql.Stmt, tx *sql.Tx) {
		for _, record := range dnsRequestBuffer {
			_, err := stmt.Exec(strings.Join(record.Answer, ","), record.DNSType, record.Duration, record.Name, record.QType, record.CreateTime)
			if err != nil {
				log.Printf("[%s]\t Failed to insert record:%v\n", logMsg, err)
				tx.Rollback()
				break
			}
		}
		dnsRequestBuffer = nil
	})
}

func syncRuleMatch(db *sql.DB) {
	logMsg, query := "RuleMatch", "INSERT INTO rule_match(duration,rule,payload,proxy,m_network,m_type,m_sourceIP,m_destinationIP,m_sourcePort,m_destinationPort,m_host,m_dnsMode,m_processPath,createTime) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	fmt.Printf("[%s]\t Writing %d records to the database...\n", logMsg, len(ruleMatchBuffer))
	startInsert(db, logMsg, query, func(stmt *sql.Stmt, tx *sql.Tx) {
		for _, record := range ruleMatchBuffer {
			_, err := stmt.Exec(record.Duration, record.Rule, record.Payload, record.Proxy, record.Metadata.Network, record.Metadata.Type, record.Metadata.SourceIP, record.Metadata.DestinationIP, record.Metadata.SourcePort, record.Metadata.DestinationPort, record.Metadata.Host, record.Metadata.DNSMode, record.Metadata.ProcessPath, record.CreateTime)
			if err != nil {
				log.Printf("[%s]\t Failed to insert record:%v\n", logMsg, err)
				tx.Rollback()
				break
			}
		}
		ruleMatchBuffer = nil
	})
}

func syncProxyDial(db *sql.DB) {
	logMsg, query := "ProxyDial", "INSERT INTO proxy_dial(address,duration,host,proxy,createTime) values(?,?,?,?,?)"
	fmt.Printf("[%s]\t Writing %d records to the database...\n", logMsg, len(proxyDialBuffer))
	startInsert(db, logMsg, query, func(stmt *sql.Stmt, tx *sql.Tx) {
		for _, record := range proxyDialBuffer {
			_, err := stmt.Exec(record.Address, record.Duration, record.Host, record.Proxy, record.CreateTime)
			if err != nil {
				log.Printf("[%s]\t Failed to insert record:%v\n", logMsg, err)
				tx.Rollback()
				break
			}
		}
		proxyDialBuffer = nil
	})
}

func syncTraffic(db *sql.DB) {
	logMsg, query := "Traffic", "INSERT INTO traffic(up,down,createTime) values(?,?,?)"
	fmt.Printf("[%s]\t Writing %d records to the database...\n", logMsg, len(trafficBuffer))
	startInsert(db, logMsg, query, func(stmt *sql.Stmt, tx *sql.Tx) {
		for _, record := range trafficBuffer {
			_, err := stmt.Exec(record.Up, record.Down, record.CreateTime)
			if err != nil {
				log.Printf("[%s]\t Failed to insert record:%v\n", logMsg, err)
				tx.Rollback()
				break
			}
		}
		trafficBuffer = nil
	})
}

func SyncDatabase(db *sql.DB) {
	mutex.Lock()
	syncDNSRequest(db)
	syncProxyDial(db)
	syncRuleMatch(db)
	syncTraffic(db)
	mutex.Unlock()
}
