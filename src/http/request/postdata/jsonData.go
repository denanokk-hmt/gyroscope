/*
======================
POSTリクエストされたパラメーターの受信
JsonDataを処理
========================
*/
package postdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	COMMON "bwing.app/src/common"
	REQ "bwing.app/src/http/request/request"
)

// /////////////////////////////////////////////////
func ParseJson(r *http.Request, rq *REQ.RequestData) ([]uint8, int, error) {

	var err error

	//Validation
	if rq.Method == "GET" {
		err = errors.New("Method")
		return nil, http.StatusBadRequest, err
	}

	//Check Contente-TYpe
	rcts := r.Header.Get("Content-Type")                           //リクエストしてきたContentType
	ctArr := strings.Split(strings.ReplaceAll(rcts, " ", ""), ";") //リクエストしてきたContentTypeを分解
	dcts := REQ.GetContentTypeAll()                                //デフォルトの許可ContentTypeを取得
	var cBool bool = false
	for _, ct := range dcts {
		if COMMON.StringSliceSearch(ctArr, ct) {
			cBool = true
			break
		}
	}
	if !cBool {
		err = errors.New("Content-Type")
		return nil, http.StatusBadRequest, err
	}

	//To allocate slice for request body
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		err = errors.New("Content-Length")
		return nil, http.StatusInternalServerError, err
	}

	//Read body data to parse json
	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		err = errors.New("")
		return nil, http.StatusInternalServerError, err
	}

	return body, length, nil
}

///////////////////////////////////////////////////
/* ===========================================
JSON形式のPOSTパラメーターをパースして格納
=========================================== */
func ParseJsonData(r *http.Request, rq *REQ.RequestData) error {

	//parse json request body
	body, length, err := ParseJson(r, rq)
	if err != nil {
		return err
	}

	//RequestDataに格納
	rq.Body = string(body)

	//parse json
	var jsonBody map[string]interface{}
	err = json.Unmarshal(body[:length], &jsonBody)
	if err != nil {
		s := fmt.Sprintf("[Gyroscope] Unmarshal error:[%s] body:[%s]", err, string(body))
		err = errors.New(s)
		return err
	}

	//fmt.Printf("%v\n", jsonBody)
	var pSs string //Params stringsの箱

	//元の型へキャストしてPrameterへ格納、連想配列が"IDKey"or"NameKey"の場合は、Keyへ格納
	var p REQ.PostParameter
	for n, v := range jsonBody {
		switch v := v.(type) {
		case string:
			//fmt.Printf("%s:%s (%T)\n", n, v, v)
			s := jsonBody[n].(string)
			p = REQ.PostParameter{Name: n, Type: "string", StringValue: s}
			pSs = pSs + n + "=" + s + "|"
		case int:
			//fmt.Printf("%s:%d (%T)\n", n, v, v)
			i := int(jsonBody[n].(int))
			p = REQ.PostParameter{Name: n, Type: "int", IntValue: i}
			pSs = pSs + n + "=" + strconv.Itoa(i) + "|"
		case int64:
			//fmt.Printf("%s:%d (%T)\n", n, v, v)
			i64 := int64(jsonBody[n].(int64))
			p = REQ.PostParameter{Name: n, Type: "int64", Int64Value: i64}
			pSs = pSs + n + "=" + strconv.FormatInt(i64, 10) + "|"
		case float32:
			//fmt.Printf("%s:%f (%T)\n", n, v, v)
			f32 := float32(jsonBody[n].(float32))
			p = REQ.PostParameter{Name: n, Type: "float32", Float32Value: f32}
			s := fmt.Sprintf("%f", f32)
			pSs = pSs + n + "=" + s + "|"
		case float64:
			//fmt.Printf("%s:%f (%T)\n", n, v, v)
			f64 := float64(jsonBody[n].(float64))
			p = REQ.PostParameter{Name: n, Type: "float64", Float64Value: f64}
			s := strconv.FormatFloat(f64, 'f', 0, 64)
			pSs = pSs + n + "=" + s + "|"
		case bool:
			//fmt.Printf("%s:%t (%T)\n", n, v, v)
			b := jsonBody[n].(bool)
			p = REQ.PostParameter{Name: n, Type: "bool", BoolValue: b}
			s := strconv.FormatBool(b)
			pSs = pSs + n + "=" + s + "|"
		//case Robot:
		//	fmt.Printf("%s: %+v (%T)\n", i, v, v) // valStruct: {name:Doraemon birth:2112} (main.Robot)
		default:
			//fmt.Println(p)
			//fmt.Printf("I don't know about type %s %T!\n", n, v)
			if v == nil {
				p = REQ.PostParameter{Name: n, Type: "string", StringValue: ""}
				pSs = pSs + n + "=" + "" + "|"
			}
		}
		//格納
		rq.PostParameter = append(rq.PostParameter, p)

		//log_nameに指定があった場合、これをロギングのjsonPayload.logNameに指定する
		if n == "log_name" {
			s := jsonBody[n].(string)
			rq.ParamsBasic.LogName = s
		}
	}

	//リクエストIPを取得し、パラメーターに格納。リクエストIPをParams文字列に追記
	ip := r.Header.Get("x-forwarded-for")
	p = REQ.PostParameter{Name: "Ip", Type: "string", StringValue: ""}
	rq.PostParameter = append(rq.PostParameter, p)
	pSs = pSs + "ip" + "=" + ip + "|"

	//timestamp
	nowDt := time.Now().UTC()
	layout := "2006-01-02T15:04:05.000000000Z"
	st := nowDt.Format(layout)
	pSs = pSs + "timestamp" + "=" + st + "|"

	//パラメーターを文字列で格納
	rq.ParamsBasic.ParamsStrings = pSs

	return nil
}
