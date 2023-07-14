/*
======================
POST Methodリクエストに対する処理を行わせ、結果をレスポンスする
========================
*/
package response

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ERR "bwing.app/src/error"
	REQ "bwing.app/src/http/request/request"

	"net/http"

	COMMON "bwing.app/src/common"
	CONFIG "bwing.app/src/config"
	FILE "bwing.app/src/file"
	LOG "bwing.app/src/log"
)

// Inerface
type PostResponse struct {
	InsertId string //ロギング番号
}

///////////////////////////////////////////////////
/* ===========================================
イベントログをトラッキング
=========================================== */
func (res PostResponse) TrackingEventsLog(w http.ResponseWriter, r *http.Request, rq *REQ.RequestData) {

	/*-----------------------------
	ログの整形とファイル書き込み(更新、作成)
	----------------------------- */

	//Postデータを文字列に整形
	eventLog, sDate := LOG.GyroscopeLogging(rq, res.InsertId, LOG.INFO)

	//イベントログの書き込み
	f := &FILE.EventsLog{Sdate: sDate}
	fName, err := f.WriteEventsLog2File(eventLog)
	if err != nil {
		ERR.ErrorResponse(w, rq, err, http.StatusInternalServerError, true)
		return
	}

	//結果
	responseOutput := fmt.Sprintf("[Gyroscope] Events log received and memorized.[%s]", fName)

	//結果を出力
	fmt.Println(LOG.SetLogEntry2("", LOG.INFO, CONFIG.LOGGING_NAME, responseOutput))

	//結果をレスポンス
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode((responseOutput))
}

///////////////////////////////////////////////////
/* ===========================================
イベントログファイル名をすべて参照
※アップロード済みの残ファイルを検索したい場合、パラメーターに、
	"search_dir":"uploaded"を指定
=========================================== */
func (res PostResponse) SearchLogStorage(w http.ResponseWriter, r *http.Request, rq *REQ.RequestData) {

	//削除するファイル名パラメータを取得
	var sDir string
	for _, p := range rq.PostParameter {
		if p.Name == "search_dir" {
			sDir = p.StringValue
			break
		}
	}

	//ファイル名を取得
	fileNames, err := FILE.GetEventsLogFileNames(sDir)
	if err != nil {
		ERR.ErrorResponse(w, rq, err, http.StatusInternalServerError, true)
		return
	}

	//結果
	responseOutput := fmt.Sprintf("%v\n", fileNames)

	//結果を出力
	fmt.Println(LOG.SetLogEntry2("", LOG.INFO, CONFIG.LOGGING_NAME, responseOutput))

	//結果をレスポンス
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode((responseOutput))
}

///////////////////////////////////////////////////
/* ===========================================
強制的に最終ログをGCSへアップロードする
※バッチを実行させる場合、パラメーターに、
	"upload_filename":"*"を指定
=========================================== */
func (res PostResponse) UploadLogStorage(w http.ResponseWriter, r *http.Request, rq *REQ.RequestData) {

	/*-----------------------------
	準備(日時、ファイル名)
	----------------------------- */

	//現在日時を設定
	nt := time.Now().UTC()
	dt := nt.Format(CONFIG.LOG_FILE_LAYOUT)

	f := &FILE.EventsLog{Sdate: dt}

	//ファイル名を確認し出力(書き込み用のログファイル
	fileNames, err := FILE.GetEventsLogFileNames("")
	if err != nil {
		ERR.ErrorResponse(w, rq, err, http.StatusInternalServerError, true)
		return
	}

	//アップロードするファイル名パラメータを取得
	var uf string
	for _, p := range rq.PostParameter {
		if p.Name == "upload_filename" {
			uf = p.StringValue
			break
		}
	}

	//アスタリスク指定の場合、実行日と比較してファイル日時が古いものすべてが対象
	if uf != "*" {
		for _, fn := range fileNames {
			if fn == uf {
				fileNames = nil
				fileNames = append(fileNames, fn)
				f.ForceUpload = true
				break
			}
		}
	}

	//Gard
	if len(fileNames) == 0 {
		err := fmt.Errorf("not found upload file")
		ERR.ErrorResponse(w, rq, err, http.StatusInternalServerError, true)
		return
	}

	/*-----------------------------
	アップロード処理
	----------------------------- */

	//アップファイル名をセット
	f.FileNames = fileNames

	//GCSアップロード
	ufs, err := f.UploadEventsLog2Gcs()
	if err != nil {
		ERR.ErrorResponse(w, rq, err, http.StatusInternalServerError, true)
		return
	}

	//結果
	responseOutput := fmt.Sprintf("Uploaded:%v", ufs)

	//結果を出力
	fmt.Println(LOG.SetLogEntry2("", LOG.INFO, CONFIG.LOGGING_NAME, responseOutput))

	//結果をレスポンス
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode((responseOutput))
}

///////////////////////////////////////////////////
/* ===========================================
イベントログファイルをストレージから削除
※バッチを実行させる場合、パラメーターに、
	"delete_dir":"uploaded"を指定
	"delete_filename":"*"を指定
=========================================== */
func (res PostResponse) DeleteLogStorage(w http.ResponseWriter, r *http.Request, rq *REQ.RequestData) {

	/*-----------------------------
	準備
	----------------------------- */

	//削除するファイル名パラメータを取得
	var dDir string
	for _, p := range rq.PostParameter {
		if p.Name == "delete_dir" {
			dDir = p.StringValue
			break
		}
	}

	//ファイル名を確認し出力(書き込み用のログファイル
	fileNames, err := FILE.GetEventsLogFileNames(dDir)
	if err != nil {
		ERR.ErrorResponse(w, rq, err, http.StatusInternalServerError, true)
		return
	}

	//削除するファイル名パラメータを取得
	var dRe string
	for _, p := range rq.PostParameter {
		if p.Name == "delete_regexp" {
			dRe = p.StringValue
			break
		}
	}

	//アスタリスク指定の場合、すべてが削除対象、正規表現指定の場合、一致するものに選定
	if dRe != "" {
		fileNames = COMMON.MatchedRegExpSliceValue(fileNames, dRe)
	}

	//Gard
	if len(fileNames) == 0 {
		err := fmt.Errorf("not found delete file")
		ERR.ErrorResponse(w, rq, err, http.StatusInternalServerError, true)
		return
	}

	/*-----------------------------
	デリート処理
	----------------------------- */

	f := &FILE.EventsLog{FileNames: fileNames}

	//アップファイル名をセット
	f.FileNames = fileNames

	//ファイル数を確認し出力
	err = f.DeleteEventsLog(dDir)
	if err != nil {
		ERR.ErrorResponse(w, rq, err, http.StatusInternalServerError, true)
		return
	}

	//結果
	responseOutput := fmt.Sprintf("Delete Events log files. %v", fileNames)

	//結果を出力
	fmt.Println(LOG.SetLogEntry2("", LOG.INFO, CONFIG.LOGGING_NAME, responseOutput))

	//結果をレスポンス
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode((responseOutput))
}

///////////////////////////////////////////////////
/* ===========================================
イベントログファイルを生成する
※Weeklyでバッチなどで先に書き込むログファイルを生成しておく
	"term": 15 or 30 or 60
	"frequecy":"weekly"を指定
		"daily": 実行日の次の日でterm毎のログファイルを生成
		"weekly": 実行日の次の日曜日から1week分のterm毎のログファイルを生成
		"monthly": 実行日の次の月1month分のterm毎のログファイルを生成
	”start_date":
		Format:
			"2023-01-01"：日付文字列で指定した場合、その日付からファイルが作成される
			"next": 実行日の次のterm日付からファイルが作成される→バッチに使用
				daily-->次の日付
				weekly-->次の日曜日の日付
				monthly-->次の月の初日
=========================================== */
func (res PostResponse) CreateEventsLogFiles(w http.ResponseWriter, r *http.Request, rq *REQ.RequestData) {

	/*-----------------------------
	準備
	----------------------------- */

	var term int
	var frequency string
	var startDate string

	//パラメーターを取得
	for _, p := range rq.PostParameter {
		switch p.Name {
		case "term":
			term = int(p.Float64Value)
		case "frequency":
			frequency = p.StringValue
		case "start_date":
			startDate = p.StringValue
		}
	}

	//startDateが"next"だった場合、frequencyから開始日を計算
	if startDate == "next" {
		startDate = strings.Split(time.Now().String(), " ")[0]
	}

	//termバリデーション
	switch term {
	case 15, 30, 60:
	default:
		err := fmt.Errorf("valid error term is not 15 or 30 or 60:[%d]", term)
		ERR.ErrorResponse(w, rq, err, http.StatusBadRequest, true)
		return
	}

	//frequencyバリデーション
	switch frequency {
	case "daily", "weekly", "monthly":
	default:
		err := fmt.Errorf("valid error term is not daily or weekly or monthly:[%s]", frequency)
		ERR.ErrorResponse(w, rq, err, http.StatusBadRequest, true)
		return
	}

	//日付バリデーション
	b := COMMON.DateFormatChecker(startDate)
	if !b {
		err := fmt.Errorf("valid error start_date is not date(string):[%s]", startDate)
		ERR.ErrorResponse(w, rq, err, http.StatusBadRequest, true)
		return
	}

	//イベントログの生成
	f := &FILE.EventsLog{}
	err := f.CreateEventsLogFilesByTerm(term, frequency, startDate)
	if err != nil {
		ERR.ErrorResponse(w, rq, err, http.StatusInternalServerError, true)
		return
	}

	//結果
	responseOutput := fmt.Sprintf("[Gyroscope] create next events log files:[Qty:%d][term:%d][frequency:%s][startDate:%s]", len(f.FileNames), term, frequency, startDate)

	//結果を出力
	fmt.Println(LOG.SetLogEntry2("", LOG.INFO, CONFIG.LOGGING_NAME, responseOutput))

	//結果をレスポンス
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode((responseOutput))
}
