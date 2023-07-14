/*
======================
Logging
========================
*/
package log

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	CONFIG "bwing.app/src/config"
	REQ "bwing.app/src/http/request/request"
)

var (
	INFO  = "INFO"
	WARN  = "WARNING"
	ERROR = "ERROR"
)

type LogEntry struct {
	InsertId    string    `json:"insertId"`    //ログ番号
	Severity    string    `json:"severity"`    //ログレベル
	LogName     string    `json:"logName"`     //ログ名
	TextPayload string    `json:"textPayload"` //ログ内容
	Timestamp   time.Time `json:"timestamp"`   //ログタイムスタンプ(UTC)
}

type OsSettingsStruct struct {
	Hostname string
}

var OsSettings OsSettingsStruct

type RequestLogs struct {
	Host    string      `json:"host"`
	Method  string      `json:"method"`
	Urlpath string      `json:"url_path"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
	Params  string      `json:"params"`
}

func init() {
	n, _ := os.Hostname()
	OsSettings.Hostname = n
}

func init() {
	log.SetPrefix("") // 接頭辞の設定
}

///////////////////////////////////////////////////
/* ===========================================
構造体をJSON形式の文字列へ変換
=========================================== */
func (l LogEntry) String() string {
	out, err := json.Marshal(l)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

// ログエントリの箱につめる
func SetLogEntry(insertId, level, logName, text string) string {
	entry := &LogEntry{
		InsertId:    insertId,
		Severity:    level,
		LogName:     logName,
		TextPayload: text,
		Timestamp:   time.Now(),
	}
	return entry.String()
}

// ログエントリの箱につめる
func SetLogEntry2(insertId, level, logName, text string) (string, string) {
	//日付オブジェ
	nowDt := time.Now().UTC()
	layout := "2006-01-02T15:04:05"
	st := nowDt.Format(layout)

	entry := &LogEntry{
		InsertId:    insertId,
		Severity:    level,
		LogName:     logName,
		TextPayload: text,
		Timestamp:   nowDt,
	}

	//文字列化
	s := entry.String()

	//文字列化で生まれてしまう、Jsonin Jsonの余分エスケープofエスケープ"\\"を排除
	reg := regexp.MustCompile(`[\\\\]`)
	s = reg.ReplaceAllString(s, "")
	return s, st
}

///////////////////////////////////////////////////
/* ===========================================
API Request Logging
=========================================== */
func ApiRequestLogging(rq *REQ.RequestData, insertId string, level string) {

	//出力項目
	var output RequestLogs = RequestLogs{
		Host:    OsSettings.Hostname,
		Method:  rq.Method,
		Urlpath: rq.Urlpath,
		Headers: rq.Header.(http.Header),
		Body:    rq.Body,
		Params:  rq.ParamsBasic.ParamsStrings,
	}

	//リクエストをロギング
	fmt.Println(SetLogEntry2(insertId, level, "ApiRequestLogging", fmt.Sprintf("%+v", output)))
}

///////////////////////////////////////////////////
/* ===========================================
リクエストされたPostデータをすべてロギングする
=========================================== */
func GyroscopeLogging(rq *REQ.RequestData, insertId, level string) (string, string) {

	//出力項目
	var output RequestLogs = RequestLogs{
		Host:    rq.Host,
		Method:  rq.Method,
		Urlpath: rq.Urlpath,
		Headers: rq.Header.(http.Header),
		Params:  rq.ParamsBasic.ParamsStrings,
	}

	//logNameを指定
	var logName string
	if rq.ParamsBasic.LogName != "" {
		logName = rq.ParamsBasic.LogName
	} else {
		logName = CONFIG.LOGGING_NAME
	}

	//Logを整形
	eventLog, sDate := SetLogEntry2(insertId, level, logName, fmt.Sprintf("%+v", output))

	return eventLog, sDate
}

///////////////////////////////////////////////////
/* ===========================================
簡易ロギング
=========================================== */
func JustLogging(text string) {
	pod := CONFIG.GetConfig(CONFIG.ENV_MY_POD_NAME)
	if pod == "" {
		pod = CONFIG.GetConfig(CONFIG.ENV_SERVER_CODE)
	}
	fmt.Printf("[%s] %s\n", pod, text)
}
