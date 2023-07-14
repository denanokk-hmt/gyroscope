/*
======================
ファイル操作関連を扱う
========================
*/
package file

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

type Files struct {
	DirPath  string
	FileName string
	Asc      bool
}

///////////////////////////////////////////////////
/*===========================================
ファイル名を取得、ソード
※ディレクトリリカーシブルではない
===========================================*/
func (f Files) GetFileNamesAndSortSimple() ([]string, error) {

	var fileNames []string

	//ディレクトリを取得
	files, err := ioutil.ReadDir(f.DirPath)
	if err != nil {
		return nil, err
	}

	//ディレクトリ配下すべてのファイルを取得
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	//ファイルをソート
	sort.Slice(fileNames, func(i, j int) bool {
		if f.Asc {
			return fileNames[i] < fileNames[j] //asc: true 昇順
		} else {
			return fileNames[i] > fileNames[j] //asc: false 降順
		}
	})

	return fileNames, nil
}

///////////////////////////////////////////////////
/*===========================================
スキャナーで読み込み、行単位で配列にして戻す
===========================================*/
func (f Files) ReadAllByScanner() ([]string, error) {

	//ファイルOpen
	orgFile, err := os.Open(f.DirPath + f.FileName)
	if err != nil {
		return nil, err
	}
	defer orgFile.Close()

	//スキャナーにバッファ
	scanner := bufio.NewScanner(orgFile)

	//一行づつ箱に詰める
	var ss []string
	for scanner.Scan() {
		ss = append(ss, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return ss, nil
}

///////////////////////////////////////////////////
/*===========================================
ファイルOpen
モード：追記、存在していなければ新規作成
===========================================*/
func (f Files) WritingOpenFileMode_WR_CREATE_APPEND_SYNC(text string) error {

	//新規作成または、追記モードでファイルをOpen
	file, err := os.OpenFile(f.DirPath+f.FileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	//ログを追記
	fmt.Fprintln(file, text)
	return nil
}

/*
===========================================
ファイルを開かずにサイズを確認
===========================================
*/
func (f Files) GetClosedFileSize() (int64, error) {
	fileinfo, err := os.Stat(f.DirPath + f.FileName)
	if err != nil {
		return 0, err
	}
	return fileinfo.Size(), nil
}

/*
===========================================
ファイルを削除
===========================================
*/
func (f Files) DeleteFile() error {
	//アップロードしたファイルは削除
	err := os.Remove(f.DirPath + f.FileName)
	if err != nil {
		return err
	}
	return nil
}

/*
===========================================
ファイルをリネーム(移動)
===========================================
*/
func (f Files) RenameFile(destinyFilePath string) error {
	err := os.Rename(f.DirPath+f.FileName, destinyFilePath)
	if err != nil {
		return err
	}
	return nil
}
