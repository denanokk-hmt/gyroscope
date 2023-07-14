/*
======================
イベントログの書き込み処理
========================
*/
package file

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	COMMON "bwing.app/src/common"
	CONFIG "bwing.app/src/config"
	GCS "bwing.app/src/gcs"
	LOG "bwing.app/src/log"
)

type EventsLog struct {
	Sdate       string   //日時の文字列(書き込み:タイムスタンプ、バッチなどの毎時のUpload:実行日時)Foamat:"2023-02-24T11:42:04"
	ForceUpload bool     //強制的にアップロードを行う
	FileNames   []string //アップロードを行うファイル名たち
}

///////////////////////////////////////////////////
/*===========================================
Tremに従いevent logファイルを生成
===========================================*/
func (e *EventsLog) CreateEventsLogFilesByTerm(term int, frequency, startDate string) error {

	//開始日、終了日を設定
	seDate := COMMON.GetNextFrequencyDates(frequency, startDate, "-")

	//EventsLog files名を生成
	fileNames, err := GetTermEventsLogFilesNames(term, seDate[0], seDate[1], "-", "txt")
	if err != nil {
		return err
	}

	//空のログファイルを生成する
	for i, fn := range fileNames {

		//Filesの準備(ファイル格納ディレクトリ, ファイル名をセット)
		f := &Files{DirPath: CONFIG.GetConfig(CONFIG.LOG_DIR_PATH), FileName: fn}

		//新規作成または、追記モードでファイルをOpen
		err := f.WritingOpenFileMode_WR_CREATE_APPEND_SYNC("")
		if err != nil {
			return err
		}
		go LOG.JustLogging(fmt.Sprintf("create log file:[%s][%d/%d]", fn, i+1, len(fileNames)))
	}
	e.FileNames = fileNames

	return nil
}

///////////////////////////////////////////////////
/*===========================================
event log を一時間毎のファイルに書き込む
===========================================*/
func (e EventsLog) WriteEventsLog2File(eventLog string) (string, error) {

	/*-----------------------------
	準備(日時、ファイル名、ディレクトリ)
	----------------------------- */

	//ログタイムスタンプを分解してファイル名を導く
	//Foamat:"2023-02-24T11:42:04"
	sDate := strings.Split(e.Sdate, "T")
	ts := strings.Split(sDate[1], ":")
	min, _ := strconv.Atoi(ts[1])

	//Termに従いファイル名を決定
	var fName string
	term, _ := strconv.Atoi(CONFIG.LOG_TERM)
	switch term {
	case 15:
		switch {
		case min < 15:
			fName = sDate[0] + "T" + ts[0] + ":00:00_" + ts[0] + ":14:59"
		case min >= 15 && min < 30:
			fName = sDate[0] + "T" + ts[0] + ":15:00_" + ts[0] + ":29:59"
		case min >= 30 && min < 45:
			fName = sDate[0] + "T" + ts[0] + ":30:00_" + ts[0] + ":44:59"
		case min >= 45 && min < 60:
			fName = sDate[0] + "T" + ts[0] + ":45:00_" + ts[0] + ":59:59"
		}
	case 30:
		switch {
		case min < 30:
			fName = sDate[0] + "T" + ts[0] + ":00:00_" + ts[0] + ":29:59"
		case min >= 30 && min < 60:
			fName = sDate[0] + "T" + ts[0] + ":30:00_" + ts[0] + ":59:59"
		}
	case 60:
		fName = sDate[0] + "T" + ts[0] + ":00:00_" + ts[0] + ":59:59"
	}

	//変換後のファイル名Format:"2023-02-24T11:00:00_11:59:59_0.txt"
	fName += CONFIG.LOG_FILE_SUFFIX_POD_NO + ".txt"

	//Filesの準備(ファイル格納ディレクトリ, ファイル名をセット)
	f := &Files{DirPath: CONFIG.GetConfig(CONFIG.LOG_DIR_PATH), FileName: fName}

	/*-----------------------------
	ログの書き込み
	-----------------------------*/

	//新規作成または、追記モードでファイルをOpen
	err := f.WritingOpenFileMode_WR_CREATE_APPEND_SYNC(eventLog)
	if err != nil {
		return "", err
	}
	go LOG.JustLogging(fmt.Sprintf("WritingOpenFileMode_WR_CREATE_APPEND_SYNC [%s][%s]", f.DirPath, f.FileName))

	return f.FileName, nil
}

///////////////////////////////////////////////////
/*===========================================
GSCへログファイルをアップロード
===========================================*/
func (e EventsLog) UploadEventsLog2Gcs() ([]string, error) {

	/*-----------------------------
	準備(日時、ファイル名、ディレクトリ)
	----------------------------- */

	//アップロードするファイル名を格納
	var uploadFiles []string
	if e.ForceUpload {
		uploadFiles = e.FileNames //強制的にファイルをアップロード
	} else {
		//実行日時と比較して、日付やHourが古い場合のファイルを検索
		//実行日時を分解して、対象ファイルを導く
		//Foamat:"2023-03-07T10:00:00_10:14:59.txt
		sTs := COMMON.ConvertStringDateToUtcTime(e.Sdate, "T", "-", ":", true)

		//アップロードするファイル名を箱に詰める
		for _, fn := range e.FileNames {
			//Foamat:"2023-03-07T10:00:00_10:14:59.txtの後方の時間を取る
			ts := strings.Split(fn, "T")[0] + "T" + strings.Split(fn, "_")[1]
			ts = strings.Split(ts, ".")[0]
			fTs := COMMON.ConvertStringDateToUtcTime(ts, "T", "-", ":", true)

			//ファイル日時と実行日時を比較して、ファイル日時が古い場合、採用
			if fTs.Before(sTs) {
				uploadFiles = append(uploadFiles, fn)
			}
		}
	}

	//Filesの準備(ファイル格納ディレクトリ)
	fDir := CONFIG.GetConfig(CONFIG.LOG_DIR_PATH)

	//ファイルリミットに応じて分割
	if CONFIG.EVENTS_LOG_FILE_DIVIDE {
		var err error
		uploadFiles, err = devideEventsLogFileByite(fDir, uploadFiles)
		if err != nil {
			return nil, err
		}
	}

	/*-----------------------------
	GCSへアップロードとローテート(移動)
	----------------------------- */

	for i, uf := range uploadFiles {

		//アップロードDirをファイル名から取得(Format:2006-1-12T00:00:00_00:59:59.txt)
		bDir := strings.ReplaceAll(strings.Split(uf, "T")[0], "-", "/")

		LOG.JustLogging(fmt.Sprintf("GCSへアップロード中:[%s][%d/%d]", uf, i+1, len(uploadFiles)))

		//GCSバケットアップロード
		g := &GCS.GcsUpload{DirPath: fDir, FileName: uf, BucketDirPath: bDir}
		err := g.UploadFiletoGcsBucket()
		if err != nil {
			return nil, err
		}

		//アップロードしたファイルを移動
		f := &Files{DirPath: fDir, FileName: uf}
		destinyFilePath := fDir + CONFIG.LOG_UPLOADED_DIR_PATH + uf
		err = f.RenameFile(destinyFilePath)
		if err != nil {
			return nil, err
		}

		LOG.JustLogging(fmt.Sprintf("GCSへアップロード完了:[%s][%d/%d]", uf, i+1, len(uploadFiles)))
	}

	return uploadFiles, nil
}

///////////////////////////////////////////////////
/*===========================================
event log を削除
===========================================*/
func (e EventsLog) DeleteEventsLog(addDir string) error {

	//Filesの準備(ファイル格納ディレクトリ)
	dir := CONFIG.GetConfig(CONFIG.LOG_DIR_PATH)
	if addDir != "" {
		dir += addDir + "/" //例uploaded削除バッチ時に指定する
	}
	f := &Files{DirPath: dir}

	//すべてのログファイルを削除
	for _, df := range e.FileNames {
		f = &Files{DirPath: dir, FileName: df}
		err := f.DeleteFile()
		if err != nil {
			return err
		}
	}

	return nil
}

///////////////////////////////////////////////////
/*===========================================
event log ファイル名を全部戻す
===========================================*/
func GetEventsLogFileNames(addDir string) ([]string, error) {

	//Filesの準備(ファイル格納ディレクトリ)
	dir := CONFIG.GetConfig(CONFIG.LOG_DIR_PATH)
	if addDir != "" {
		dir += addDir //例uploaded削除バッチ時に指定する
	}
	f := &Files{DirPath: dir}

	//ファイル名を降順で取得
	f.Asc = false
	fileNames, err := f.GetFileNamesAndSortSimple()
	if err != nil {
		return nil, err
	}

	//NFS対応(os.Readの検索で"lost+found"がファイル名に入ってくる)、余計なファイル名を削除
	fileNames = COMMON.RemoveSliceValue(fileNames, "lost+found")

	return fileNames, nil
}

///////////////////////////////////////////////////
/*===========================================
指定期間におけるTerm毎ログファイル名を取得
	startDate, endDate: "2023-01-01"
	delimiter: "-"
	extention: "txt"
	term: 15 or 30 or 60
	Terｍが15の場合に欲しいファイル名のFormat
		2023-01-01T00:00:00_00:14:59.txt
		2023-01-01T00:15:00_00:29:59.txt
		2023-01-01T00:30:00_00:44:59.txt
		2023-01-01T00:45:00_00:59:59.txt
===========================================*/
func GetTermEventsLogFilesNames(term int, startDate, endDate, dateDelimiter, extension string) ([]string, error) {

	//指定期間の日数
	diff := COMMON.DateDiffCalculator(startDate, endDate, dateDelimiter)

	//日数をもとに期間の年月日を配列で取得
	sDates := COMMON.DateAddCalculator(startDate, "-", diff)

	var fileNames []string

	//日付期間中において、termで分割された24時間分のファイル名を配列で取得
	for _, s := range sDates {
		//24時間分回す
		for h := 0; h < 24; h++ {

			//term別の分割回数(=Loop回数)
			var termLoop int
			switch term {
			case 15:
				termLoop = 4
			case 30:
				termLoop = 2
			case 60:
				termLoop = 1
			}

			//ファイル名を生成
			var termBuff int
			var fileName string
			for i := 0; i < termLoop; i++ {
				minS := fmt.Sprintf("%02d", termBuff)                            //termバッファ=開始時刻:分
				termBuff = termBuff + term                                       //termをバッファ(累積)
				minE := fmt.Sprintf("%02d", termBuff-1)                          //termバッファから1を引く=終了時刻:分
				startT := fmt.Sprintf("%02d", h) + ":" + minS + ":00"            //開始時刻を形成
				endT := fmt.Sprintf("%02d", h) + ":" + minE + ":59." + extension //終了時刻を形成
				fileName = s + "T" + startT + "_" + endT                         //日付と時刻をつなぎ合わせる
				fileNames = append(fileNames, fileName)                          //ファイル名を格納
			}
		}
	}

	return fileNames, nil
}

///////////////////////////////////////////////////
/*===========================================
event log を分割する
===========================================*/
func devideEventsLogFileByite(fDir string, uploadFiles []string) ([]string, error) {

	var ufs []string //分割ファイル名を入れる箱

	//ファイルリミットに応じて分割
	for _, uf := range uploadFiles {

		//ソースFilesの準備(ファイル格納ディレクトリ)
		f := &Files{DirPath: fDir, FileName: uf}

		//ファイルサイズを取得
		fSize, err := f.GetClosedFileSize()
		if err != nil {
			return nil, err
		}

		//書き込みを行ったファイルサイズが設定しきい値を越えた場合、アップロード対象とする
		if fSize <= CONFIG.EVENTS_LOG_FILE_SIZE_THRESHOLD {
			ufs = uploadFiles
		} else {

			//元ファイルを読み込み
			ss, err := f.ReadAllByScanner()
			if err != nil {
				return nil, err
			}

			//チャンク数を計算
			chunks := COMMON.ChunkCalculator2(len(ss), CONFIG.UPLOAD_FILE_MAX_ROWS)

			//データのチャンク数に応じて分割ファイルを生成する
			for i, c := range chunks.Positions {

				//新しいファイル名を設定(拡張子の前に連番_S[i]を挿入)
				suffix := "_S" + strconv.Itoa(i)
				nf := strings.Split(uf, ".")[0] + suffix + ".txt"

				LOG.JustLogging(fmt.Sprintf("アップロードファイルを分割中:[%s][%d/%d]", nf, i+1, len(chunks.Positions)))

				//新規作成または、追記モードでファイルをOpen
				newFile, err := os.OpenFile(f.DirPath+f.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0666)
				if err != nil {
					return nil, err
				}
				defer newFile.Close()

				//新規作成ファイルに行単位で書き込む
				for i := c.Start; i < c.End; i++ {
					fmt.Fprintln(newFile, ss[i]) //ログを追記
				}

				//新規作成ファイル名をつめる
				ufs = append(ufs, nf)

				LOG.JustLogging(fmt.Sprintf("ファイル分割完了:[%s][%d/%d]", nf, i+1, len(chunks.Positions)))
			}

			//分割した元ファイルを移動
			f := &Files{DirPath: fDir, FileName: uf}
			destinyFilePath := fDir + CONFIG.LOG_UPLOADED_DIR_PATH + uf
			err = f.RenameFile(destinyFilePath)
			if err != nil {
				return nil, err
			}
		}
	}

	return ufs, nil
}
