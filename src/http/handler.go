package http

import (
	"net/http"
	"strings"
)

///////////////////////////////////////////////////
/* ===========================================
http.HandleFuncの引数に食わせて、ミドルウェ的にCorsの設定を注入させる
=========================================== */
func Handle(handlers ...func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		//CORS許可ヘッダーの設定
		CorsMiddleware(w)

		//ヘッダーにAuthorizationが含まれていた場合はpreflight成功
		if r.Method == "OPTIONS" {
			s := r.Header.Get("Access-Control-Request-Headers")
			if strings.Contains(s, "authorization") {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
			return
		}

		//リクエストを受ける
		for _, handler := range handlers {
			if err := handler(w, r); err != nil {
				return
			}
		}
	}
}

///////////////////////////////////////////////////
/* ===========================================
CORSの全許可設定
=========================================== */
func CorsMiddleware(w http.ResponseWriter) error {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	return nil
}
