/*
======================
リクエスト関連の処理の共通処理
========================
*/
package request

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	CONFIG "bwing.app/src/config"
)

var allowContentTypes []string //許可ContentType

var (
	ACTION_PATH_POSITION    = 3
	SECONDARY_PATH_POSITION = 4
)

// Request interface
type Requests interface {
	PostWithJsonSwitch(w http.ResponseWriter, r *http.Request) error
}

// httpリクエスト情報を格納する箱
type RequestData struct {
	Method        string
	Host          string
	Urlpath       string
	UrlPathArry   []string
	Uri           string
	Header        interface{}
	GetParameter  []GetParameter
	PostParameter []PostParameter
	ParamsBasic   ParamsBasic
	Body          string
}

// APIに応じた基本パラメーター
type ParamsBasic struct {
	LogName       string
	Action        string
	ClientId      string
	Token         string
	Kind          string
	Seconddary    string
	LastPath      string
	CurrentUrl    string
	CurrentParams []GetParameter
	Cdt           time.Time
	ParamsStrings string
}

// GETパラメーター情報を格納する箱
type GetParameter struct {
	Name  string
	Value string
}

// POSTパラメーター情報を格納する箱
type PostParameter struct {
	Name             string
	Type             string
	StringValue      string
	IntValue         int
	Int32Value       int32
	Int64Value       int64
	Float32Value     float32
	Float64Value     float64
	BoolValue        bool
	TimeValue        time.Time
	StringArray      []StringValue //文字配列 ["A", "B"]
	StringArrayArray [][]string
	//Int32Array   []int32
	//IntArray     []int //整数配列 [1, 2]
	//Int64Array []int64
	//Float32Array []float32
	//Float64Array []float64
	IDKeyValue   int64        //datastoreのIDKeyの値
	IDKeyArray   []Int64Value //datastoreのIDKeyのint64配列
	NameKeyValue string       //datastoreのNameKeyの値
}

// Array向け//Words(tag)向け等
type StringValue struct {
	Value string
}

// Array向け//__Key__(id)向け等
type Int64Value struct {
	Value int64
}

///////////////////////////////////////////////////
/* ===========================================
Initialize
* =========================================== */
func init() {
	allowContentTypes = []string{"application/json", "text/palin", ""}
}

///////////////////////////////////////////////////
/* ===========================================
Request情報を取得
* =========================================== */
func NewRequestData(r *http.Request) RequestData {
	var rq RequestData = RequestData{
		Method:      r.Method,
		Host:        r.Host,
		Urlpath:     r.URL.Path,
		UrlPathArry: strings.Split(r.URL.Path, "/"),
		Uri:         r.RequestURI,
		Header:      r.Header}
	return rq
}

///////////////////////////////////////////////////
/* ===========================================
Get ACTION名を格納する
※前提::Pathの前方から3つ目(はがAction名である
* =========================================== */
func NewActionName(r *http.Request, rq *RequestData) {
	//(UrlPathの前から3番目:/hmt/attachment/[Action]/)
	p := rq.UrlPathArry
	a := p[ACTION_PATH_POSITION]
	rq.ParamsBasic.Action = a
}

///////////////////////////////////////////////////
/* ===========================================
Get ACTION名の次を格納する
※前提::Pathの前方から3つ目(はがAction名である
* =========================================== */
func NewSecondaryName(r *http.Request, rq *RequestData) {
	//(UrlPathの前から4番目:/hmt/attachment/[Action]/[2ndary]/)
	p := rq.UrlPathArry
	a := p[SECONDARY_PATH_POSITION]
	rq.ParamsBasic.Seconddary = a
}

// /////////////////////////////////////////////////
// Get Last name
func NewLastName(r *http.Request, rq *RequestData) {
	//(Urlの後ろから1番目:/hmt/attachment/[Action]/[last])
	//Get action name from url path (exp: /hmt/attachment/[action]/[last])
	p := rq.UrlPathArry
	lp := len(p)
	l := p[lp-1]
	rq.ParamsBasic.LastPath = l
}

///////////////////////////////////////////////////
/* ===========================================
Post parameterの中でKeyを指定しているデータに対して
型をDatastoreに合わせて変換する
※共通：個別対応なし
* =========================================== */
func ConvertTypeKeyParameter(rq *RequestData) {
	//DatastoreのPropertyに合わせてCast
	for n, v := range rq.PostParameter {
		switch v.Name {
		case "IDKey":
			switch v.Type {
			case "string":
				ov, _ := strconv.ParseInt(rq.PostParameter[n].StringValue, 10, 64)
				rq.PostParameter[n].StringValue = ""
				rq.PostParameter[n].IDKeyValue = ov
			case "int":
				ov := rq.PostParameter[n].IntValue
				rq.PostParameter[n].IntValue = 0
				rq.PostParameter[n].IDKeyValue = int64(ov)
			case "float32":
				ov := rq.PostParameter[n].Float32Value
				rq.PostParameter[n].Float32Value = 0
				rq.PostParameter[n].IDKeyValue = int64(ov)
			case "float64":
				ov := rq.PostParameter[n].Float64Value
				rq.PostParameter[n].Float64Value = 0
				rq.PostParameter[n].IDKeyValue = int64(ov)
			default:
			}
			rq.PostParameter[n].Type = "IDKey"
		case "NameKey":
			ov := rq.PostParameter[n].StringValue
			rq.PostParameter[n].StringValue = ""
			rq.PostParameter[n].NameKeyValue = ov
		default:
			//fmt.Println(v)
		}
		//fmt.Println(n, v)
	}
}

///////////////////////////////////////////////////
/* ===========================================
Post parameterの型をDatastoreに合わせて変換する
※個別にproperty別で追加していく
* =========================================== */
func ConvertTypeParams(rq *RequestData) {

	var propPublishedAt string

	//DatastoreのPropertyに合わせてCast
	//*Article kindの場合、Number以外は、文字列で入ってくるので指定なし
	for n, v := range rq.PostParameter {
		switch v.Name {
		case "Number": //Number propertyをint型にする
			switch v.Type {
			case "string":
				ov, _ := strconv.Atoi(rq.PostParameter[n].StringValue)
				rq.PostParameter[n].StringValue = ""
				rq.PostParameter[n].IntValue = ov
			case "float32":
				ov := rq.PostParameter[n].Float32Value
				rq.PostParameter[n].Float32Value = 0
				rq.PostParameter[n].IntValue = int(ov)
			case "float64":
				ov := rq.PostParameter[n].Float64Value
				rq.PostParameter[n].Float64Value = 0
				rq.PostParameter[n].IntValue = int(ov)
			default:
			}
			rq.PostParameter[n].Type = "int"
		case "PublishedAt":
			propPublishedAt = v.Name
			ov := rq.PostParameter[n].TimeValue
			rq.PostParameter[n].StringValue = ""
			rq.PostParameter[n].TimeValue = ov
		default:
			//fmt.Println(v)
		}
	}

	//PublishedAtの指定がない→現在時刻
	if propPublishedAt == "" {
		p := PostParameter{Name: "PublishedAt", Type: "time.Time", TimeValue: time.Now()}
		rq.PostParameter = append(rq.PostParameter, p)
	}
}

///////////////////////////////////////////////////
/* ===========================================
Token認証を行う
	1.headerにTokenを指定した場合
		Authorization: Bearer <token>
	2.POST dataにTokenを指定した場合
		auth_token: <token>
	1. or 2. どちらか一方にtokenを指定
* =========================================== */
func AuthorizationToken(r *http.Request, rq *RequestData) bool {

	var token string
	token = r.Header.Get("Authorization")
	if token == "" {
		if rq.Method == "Get" {
			for _, v := range rq.GetParameter {
				if v.Name == "auth_token" {
					token = v.Value
				}
			}
		} else {
			for _, v := range rq.PostParameter {
				if v.Name == "AuthToken" {
					token = v.StringValue
				}
			}
		}
	}

	if token == "" {
		return false
	}

	//Token
	t := CONFIG.GetBearerToken()

	//一致指定売ればTrueを返す
	return "Bearer "+t == token
}

///////////////////////////////////////////////////
/* =================================
	//ContentTypeの返却
* ================================= */
func GetContentTypeAll() []string {
	return allowContentTypes
}
