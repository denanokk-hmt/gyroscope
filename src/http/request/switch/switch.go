/*
======================
リクエストのインターフェース
========================
*/
package request

import (
	"net/http"

	POST "bwing.app/src/http/request/post"
	REQ "bwing.app/src/http/request/request"
)

type request struct {
}

func NewRequests() REQ.Requests {
	return &request{}
}

///////////////////////////////////////////////////
/* ===========================================
//Interface reciver for datastore
//一枚のファイルだと縦にみづらいので、メソッド別に分離
=========================================== */
func (ra *request) PostWithJsonSwitch(w http.ResponseWriter, r *http.Request) error {
	rq := REQ.NewRequestData(r)
	POST.PostWithJson(w, r, &rq)
	return nil
}
