package utils

import (
	"database/sql"
	"embed"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type RouterItem map[string]func(http.ResponseWriter, *http.Request, *sql.DB, embed.FS)

func GetRouters() []RouterItem {
	apiList := []RouterItem{
		{"/": HandleStatic},
		{"/sync": HandleSyncDatabase},
		{"/dnsRequest": HandleDNSRequestApi},
		{"/ruleMatch": HandleRuleMatchApi},
		{"/proxyDial": HandleProxyDialApi},
		{"/traffic": HandleTrafficApi},
		{"/processDetail": HandleProcessDetailApi},
	}
	return apiList
}

func HandleDNSRequestApi(w http.ResponseWriter, r *http.Request, db *sql.DB, f embed.FS) {
	dnsRequest := QueryDnsRequests(db)
	b, _ := json.Marshal(dnsRequest)
	w.Write(b)
}

func HandleRuleMatchApi(w http.ResponseWriter, r *http.Request, db *sql.DB, f embed.FS) {
	ruleMatch := QueryRuleMatchRequests(db)
	b, _ := json.Marshal(ruleMatch)
	w.Write(b)
}

func HandleProxyDialApi(w http.ResponseWriter, r *http.Request, db *sql.DB, f embed.FS) {
	proxyDial := QueryProxyDialRequests(db)
	b, _ := json.Marshal(proxyDial)
	w.Write(b)
}

func HandleTrafficApi(w http.ResponseWriter, r *http.Request, db *sql.DB, f embed.FS) {
	traffic := QueryTrafficRequests(db)
	b, _ := json.Marshal(traffic)
	w.Write(b)
}

func HandleProcessDetailApi(w http.ResponseWriter, r *http.Request, db *sql.DB, f embed.FS) {
	query := r.URL.Query()
	path := query.Get("path")
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	processDetail := QueryProcessDetail(db, path, page, pageSize)
	b, _ := json.Marshal(processDetail)
	w.Write(b)
}

func HandleSyncDatabase(w http.ResponseWriter, r *http.Request, db *sql.DB, f embed.FS) {
	SyncDatabase(db)
	w.WriteHeader(http.StatusOK)
	w.Header().Clone()
}

func HandleStatic(w http.ResponseWriter, r *http.Request, db *sql.DB, f embed.FS) {
	url, path := "static", r.URL.Path
	if path == "/" {
		url += "/index.html"
	} else {
		url += path
	}
	b, err := f.ReadFile(url)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	arr := strings.Split(url, ".")
	suffix := arr[len(arr)-1]
	header := w.Header()
	switch suffix {
	case "css":
		header.Set("Content-Type", "text/css; charset=utf-8")
	case "js":
		header.Set("Content-Type", "text/javascript; charset=utf-8")
	case "html":
		header.Set("Content-Type", "text/html; charset=utf-8")
	}
	w.Write(b)
}
