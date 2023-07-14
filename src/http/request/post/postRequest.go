/*
======================
Postリクエスト処理
========================
*/
package request

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	COMMON "bwing.app/src/common"
	ERR "bwing.app/src/error"
	POSTD "bwing.app/src/http/request/postdata"
	REQ "bwing.app/src/http/request/request"
	RES "bwing.app/src/http/response"
	LOG "bwing.app/src/log"

	"github.com/pkg/errors"
)

///////////////////////////////////////////////////
/* ===========================================
//POST  Method request
=========================================== */
func PostWithJson(w http.ResponseWriter, r *http.Request, rq *REQ.RequestData) error {

	var err error

	/*--------------------------------
	準備
	--------------------------------*/
	cdt := time.Now() //処理開始時間
	rq.ParamsBasic.Cdt = cdt
	REQ.NewActionName(r, rq)    //URLPathからDSのAction名を取得
	REQ.NewSecondaryName(r, rq) //URLPathからAction名の次を取得

	//Jsonデータを取得
	err = POSTD.ParseJsonData(r, rq)
	if err != nil {
		ERR.ErrorResponse(w, rq, err, http.StatusBadRequest, true)
		return nil
	}

	//Paramsの型を変換
	REQ.ConvertTypeParams(rq)

	//ロギングinsertIdを設定
	insertId := COMMON.RandomString(16)

	//リクエストをロギングする
	go LOG.ApiRequestLogging(rq, insertId, LOG.INFO)

	//BearerToken認証
	b := REQ.AuthorizationToken(r, rq)
	if !b {
		err = errors.New("Authentication error.")
		ERR.ErrorResponse(w, rq, err, http.StatusUnauthorized, false)
		return nil
	}

	//Instance interface
	var response RES.PostResponse

	/***Entityを登録する各処理へのキッカー***/
	//末尾に指定されたKind名がNsKindsで登録されていない場合、[NoDs]が付与されて以下のCaseで除外される
	switch 0 {
	/* =================================================================================
	■イベントをトラッキング
	クライアントで発生したイベント情報をロギングし、
	GCSバケットにアップロードする
	*/
	case strings.Index(rq.Urlpath, "/hmt/post/tracking/events"):
		response.InsertId = insertId
		response.TrackingEventsLog(w, r, rq)

	/* =================================================================================
	■ストレージにあるイベントログファイル名を参照
	*/
	case strings.Index(rq.Urlpath, "/hmt/post/search/logstorage/events"):
		response.SearchLogStorage(w, r, rq)

		/* =================================================================================
		■ストレージにあるイベントログファイル名を強制機にGCSへアップロード
		*/
	case strings.Index(rq.Urlpath, "/hmt/post/upload/logstorage/events"):
		response.UploadLogStorage(w, r, rq)

	/* =================================================================================
	■ストレージにあるイベントログファイルを削除
	*/
	case strings.Index(rq.Urlpath, "/hmt/post/delete/logstorage/events"):
		response.DeleteLogStorage(w, r, rq)
		/* =================================================================================
		■ストレージにログファイルを先に生成しておく
		*/
	case strings.Index(rq.Urlpath, "/hmt/post/create/logstorage/events/files"):
		response.CreateEventsLogFiles(w, r, rq)
	/* =================================================================================
	例外処理
	*/
	default:
		err := errors.New("【Error】path is nothing.")
		if err != nil {
			ERR.ErrorResponse(w, rq, err, http.StatusBadRequest, false)
		}
	}

	//処理時間計測を出力
	LOG.JustLogging(fmt.Sprintf("Finish!! API[%s] 処理時間(sec): %vns", rq.Urlpath, time.Since(cdt).Seconds()))

	return nil
}
