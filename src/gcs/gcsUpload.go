/*
======================
GCSファイルアップロード
========================
*/
package gcs

import (
	"context"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"

	CONFIG "bwing.app/src/config"
)

type GcsUpload struct {
	DirPath       string //Uploadするファイルのディレクトリ
	FileName      string //Uploadするファイル名
	BucketDirPath string //GCSのパス移行のディレクトリ
}

///////////////////////////////////////////////////
/* ===========================================
ログファイルをGCSバケットにアップロード
=========================================== */
func (g GcsUpload) UploadFiletoGcsBucket() error {

	//ログを書き込まれたファイルを取得
	f, err := os.Open(g.DirPath + g.FileName)
	if err != nil {
		return err
	}
	defer f.Close()

	//Contextを設定
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	//GCS clientを生成
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	//アップロード先を設定
	bucketName := CONFIG.GetConfig(CONFIG.ENV_GCS_BUCKET_NAME)
	bucketPath := CONFIG.GetConfig(CONFIG.ENV_GCS_BUCKET_PATH)
	objectPath := bucketPath + "/" + g.BucketDirPath + "/" + g.FileName

	//オブジェクトのWriterを作成し、GCSバケットにファイルをアップロード
	writer := client.Bucket(bucketName).Object(objectPath).NewWriter(ctx)
	if _, err := io.Copy(writer, f); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}
