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
	bufferSize       = 256
	interval         = 256
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
		name, duration, count := "", 0.0, 0
		rows.Scan(&name, &duration, &count)
		dnsRequest.Counting = append(dnsRequest.Counting, tracingstruct.DNSCounting{
			Name:     name,
			Duration: duration,
			Count:    count,
		})
	}

	// DNS queryType 统计
	rows, _ = db.Query(`SELECT qType, AVG(duration) as duration, COUNT(1) as count FROM "dns_request" GROUP BY qType ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		qType, duration, count := "", 0.0, 0
		rows.Scan(&qType, &duration, &count)
		dnsRequest.QTypeCounting = append(dnsRequest.QTypeCounting, tracingstruct.DNSQTypeCounting{
			QType:    qType,
			Duration: duration,
			Count:    count,
		})
	}

	// DNS dnsType统计
	rows, _ = db.Query(`SELECT dnsType, AVG(duration) as duration, COUNT(1) as count FROM "dns_request" GROUP BY dnsType ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		dnsType, duration, count := "", 0.0, 0
		rows.Scan(&dnsType, &duration, &count)
		dnsRequest.DNSTypeCounting = append(dnsRequest.DNSTypeCounting, tracingstruct.DNSTypeCounting{
			DNSType:  dnsType,
			Duration: duration,
			Count:    count,
		})
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
		duration, payload, proxy, count := 0.0, "", "", 0
		rows.Scan(&duration, &payload, &proxy, &count)
		ruleMatch.Counting = append(ruleMatch.Counting, tracingstruct.RuleMatchCounting{
			Duration: duration,
			Payload:  payload,
			Proxy:    proxy,
			Count:    count,
		})
	}

	// 查询端口来匹配流量的类型
	rows, _ = db.Query(`SELECT m_destinationPort, COUNT(1) as count from rule_match GROUP BY m_destinationPort ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		port, count := 0, 0
		rows.Scan(&port, &count)
		ruleMatch.PortCounting = append(ruleMatch.PortCounting, tracingstruct.RuleMatchPortCounting{
			Port:  port,
			Count: count,
		})
	}

	// 查询路径来统计各个APP的连接数
	rows, _ = db.Query(`SELECT m_processPath, COUNT(1) as count FROM "rule_match" GROUP BY m_processPath ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		path, count := "", 0
		rows.Scan(&path, &count)
		ruleMatch.ProcessCounting = append(ruleMatch.ProcessCounting, tracingstruct.RuleMatchProcessCounting{
			Path:  path,
			Count: count,
		})
	}

	// 统计各个客户端IP的请求次数
	rows, _ = db.Query(`SELECT m_sourceIP, COUNT(1) as count FROM "rule_match" GROUP BY m_sourceIP ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		ip, count := "", 0
		rows.Scan(&ip, &count)
		ruleMatch.ClientCounting = append(ruleMatch.ClientCounting, tracingstruct.RuleMatchClientCounting{
			IP:    ip,
			Count: count,
		})
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
		proxy, duration, count := "", 0.0, 0
		rows.Scan(&proxy, &duration, &count)
		proxyDial.ProxyCounting = append(proxyDial.ProxyCounting, tracingstruct.ProxyDialProxyCounting{
			Proxy:    proxy,
			Duration: duration,
			Count:    count,
		})
	}

	// 查询域名被访问次数和时间
	rows, _ = db.Query(`SELECT address, AVG(duration) as duration, COUNT(1) as count FROM proxy_dial GROUP BY address ORDER BY count DESC`)
	defer rows.Close()

	for rows.Next() {
		host, duration, count := "", 0.0, 0
		rows.Scan(&host, &duration, &count)
		proxyDial.HostCounting = append(proxyDial.HostCounting, tracingstruct.ProxyDialHostCounting{
			Host:     host,
			Duration: duration,
			Count:    count,
		})
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
	rows, _ := db.Query(`SELECT AVG(up) as up, AVG(down) as down FROM "traffic"`)
	defer rows.Close()

	for rows.Next() {
		up, down := 0.0, 0.0
		rows.Scan(&up, &down)
		traffic.Up = up
		traffic.Down = down
	}

	// 查询流量历史
	rows, _ = db.Query(`SELECT up, down, createTime FROM "traffic" ORDER BY createTime DESC LIMIT 10`)
	defer rows.Close()

	for rows.Next() {
		up, down, createTime := 0.0, 0.0, 0
		rows.Scan(&up, &down, &createTime)
		traffic.History = append(traffic.History, tracingstruct.ApiTrafficHistory{
			Up:         up,
			Down:       down,
			CreateTime: createTime,
		})
	}

	mutex.Unlock()

	return traffic
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
