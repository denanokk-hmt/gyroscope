/*
=================================
サーバー起動時の起点ファイル
HttpサーバーのListen、各種初期設定を行うスタンバイ
=================================
*/
package main

import (
	"fmt"

	CONFIG "bwing.app/src/config"
	HTTP "bwing.app/src/http"
	LOG "bwing.app/src/log"
)

func main() {

	//認証用のTokenを格納
	CONFIG.NewBearerToken()

	//Get Env Configs
	configMapEnv := CONFIG.GetConfigAll()

	//Output message finish config settings.
	go LOG.JustLogging(fmt.Sprintf("[Project:%s][ServerCode:%s][Appli:%s][Env:%s][Uuid:%s]",
		configMapEnv[CONFIG.ENV_GCP_PROJECT_ID],
		configMapEnv[CONFIG.ENV_SERVER_CODE],
		configMapEnv[CONFIG.ENV_APPLI_NAME],
		configMapEnv[CONFIG.ENV_ENV],
		configMapEnv[CONFIG.ENV_MY_UUID]))

	//Output componet verion
	go LOG.JustLogging(fmt.Sprintf("Gyroscope component version :%s", CONFIG.SYSTEM_COMPONENT_VERSION))

	//Request routing & Listen PORT
	HTTP.HandleRequests()
}
