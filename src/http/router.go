/*
======================
HTTPサーバーのRouting
========================
*/
package http

import (
	"fmt"
	"log"
	"net/http"

	RC "bwing.app/src/http/request/switch"
	LOG "bwing.app/src/log"
)

///////////////////////////////////////////////////
/* ===========================================
HTTP Request Router
=========================================== */
func HandleRequests() {

	//Root Routing for inner health check
	http.HandleFunc("/", Handle(root))

	//Set Request Method Interface::厳密なRestAPIの方式には従わない(基本的にPOSTメソッドで登録・取得を行なっている)
	rc := RC.NewRequests()

	//パスの先頭は、"hmt"、その次にMethodを指定
	//Methodで処理をスイッチさせるようにswithcさせている
	http.HandleFunc("/hmt/post/", Handle(rc.PostWithJsonSwitch))

	//Listen
	port := "8080"
	LOG.JustLogging(fmt.Sprintf("GoGo Gyroscope!!%s", port))

	log.Fatal(http.ListenAndServe(":"+port, nil))

}

///////////////////////////////////////////////////
/* ===========================================
Root response for INNER Health check
=========================================== */
func root(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "HELLO ROOT!!")
	LOG.JustLogging("Endpoint Hit: root")
	return nil
}
