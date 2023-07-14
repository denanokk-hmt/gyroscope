/*
=================================
サーバーのConfigを設定する
=================================
*/
package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	COMMON "bwing.app/src/common"
)

var (
	SYSTEM_COMPONENT_VERSION = "1.2.x_gy" //Component version

	//基本環境値
	ENV_LOCAL_OR_REMORE = "LocalOrRemote"
	ENV_GCP_PROJECT_ID  = "GcpProjectId"
	ENV_SERVER_CODE     = "ServerCode"
	ENV_APPLI_NAME      = "AppliName"
	ENV_ENV             = "Env"
	ENV_GCS_BUCKET_NAME = "GcsBucketName"
	ENV_GCS_BUCKET_PATH = "GcsBucketPath"
	ENV_MY_POD_NAME     = "MyPodName"
	ENV_MY_UUID         = "MyUuid"

	LOGGING_NAME = "GyroscopeEventsLogging"

	//デフォルト
	DEFAULT_STRING_VALUE_DAILY   = "daily"
	DEFAULT_STRING_VALUE_MONTHLY = "monthly"
	DEFAULT_STRING_VALUE_TRUE    = "true"
	DEFAULT_STRING_VALUE_FALSE   = "false"
	DEFAULT_STRING_VALUE_15      = "15"
	DEFAULT_STRING_VALUE_30      = "30"
	DEFAULT_STRING_VALUE_60      = "60"

	//GCS UPLOAD FILE関連
	_, b, _, _                     = runtime.Caller(0)
	root                           = filepath.Join(filepath.Dir(b), "../../")
	LOG_STORAGE_DIR_ABSOLUTE_PATH  = root + "/cmd/logStorage/events/" //Rootディレクトリを取得して、絶対パスを指定
	LOG_STORAGE_VOLUME_DIR_PATH    = "/mnt/gyroscope/events/"         //NFSマウントしたパス
	LOG_FILE_LAYOUT                = "2006-01-02T15:04:05"
	EVENTS_LOG_FILE_DIVIDE         = false //EVENTS_LOG_FILE_SIZE_THRESHOLDを超えるサイズだった場合、UPLOAD_FILE_MAX_ROWSで分割
	EVENTS_LOG_FILE_SIZE_THRESHOLD = int64(102400000)
	LOG_DIR_PATH                   = ""
	UPLOAD_FILE_MAX_ROWS           = 100000
	LOG_UPLOADED_DIR_PATH          = "uploaded/"
	LOG_TERM                       = DEFAULT_STRING_VALUE_15
	LOG_FILE_SUFFIX_POD_NO         = "" //k8s Statefulsetにより割り振られるNo(Depoymentの場合、一意な文字列となる=logging-jobsを利用しない場合)
)

var configMapEnv map[string]string //サーバーコンフィグ値の箱
var bearerToken string             //認証Token

// /////////////////////////////////////////////////
// 起動時にGCP ProjectID、NS, Kindを登録する
func init() {

	//Set environ values
	NewConfigEnv()
}

///////////////////////////////////////////////////
/* =================================
	環境変数の格納
		$PORT
		$GCP_PROJECT_ID
		$SERVER_CODE
		$APPLI_NAME
		$ENV
* ================================= */
func NewConfigEnv() {

	//環境変数をMapping
	configMapEnv = make(map[string]string)
	configMapEnv[ENV_LOCAL_OR_REMORE] = os.Getenv("LOCAL_OR_REMOTE")
	configMapEnv[ENV_GCP_PROJECT_ID] = os.Getenv("GCP_PROJECT_ID")
	configMapEnv[ENV_SERVER_CODE] = os.Getenv("SERVER_CODE")
	configMapEnv[ENV_APPLI_NAME] = os.Getenv("APPLI_NAME")
	configMapEnv[ENV_ENV] = os.Getenv("ENV")
	configMapEnv[ENV_GCS_BUCKET_NAME] = os.Getenv("GCS_BUCKET_NAME")
	configMapEnv[ENV_GCS_BUCKET_PATH] = os.Getenv("GCS_BUCKET_PATH")

	//mount先のNFSパス
	if configMapEnv[ENV_LOCAL_OR_REMORE] == "local" {
		configMapEnv[LOG_DIR_PATH] = LOG_STORAGE_DIR_ABSOLUTE_PATH
	} else {
		configMapEnv[LOG_DIR_PATH] = LOG_STORAGE_VOLUME_DIR_PATH
	}

	//PodName
	pn := os.Getenv("MY_POD_NAME")
	configMapEnv[ENV_MY_POD_NAME] = pn

	//ログファイルのサフィックスをPodNameの末尾の文字列で指定
	//例：
	//StatefulSet:(ベース)　svc-hmt-gyroscope-0,1,2...
	//Deployment:(logging-jobsを利用せずに、自身で処理を完結させる場合) svc-hmt-gyroscope-abcsifz-04gsu,,
	pns := strings.Split(pn, "-")
	LOG_FILE_SUFFIX_POD_NO = "_" + pns[len(pns)-1]

	//現在未使用
	uuid, err := COMMON.CreateUuidV4()
	if err != nil {
		log.Fatal(err)
	}
	configMapEnv[ENV_MY_UUID] = uuid
}

///////////////////////////////////////////////////
/* =================================
	//Configの返却
* ================================= */
func GetConfig(name string) string {
	return configMapEnv[name]
}
func GetConfigAll() map[string]string {
	return configMapEnv
}

///////////////////////////////////////////////////
/* =================================
認証に用いるトークンをJSONファイルから取得しておく
* ================================= */
func NewBearerToken() {

	//箱を準備
	type BearerTokenJson struct {
		BearerTokens string `json:"bearer_token"`
	}

	//Rootディレクトリを取得して、tokensのJSONファイルの絶対パスを指定
	var (
		_, b, _, _ = runtime.Caller(0)
		root       = filepath.Join(filepath.Dir(b), "../../")
	)
	path := root + "/authorization/bearer_token.json"

	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// JSONデコード
	var token BearerTokenJson
	if err := json.Unmarshal(bytes, &token); err != nil {
		log.Fatal(err)
	}
	bearerToken = token.BearerTokens

}

///////////////////////////////////////////////////
/* =================================
認証に用いるトークンを取得
* ================================= */
func GetBearerToken() string {
	return bearerToken
}
