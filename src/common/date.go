/*
======================
共通処理のファイル
========================
*/
package common

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

///////////////////////////////////////////////////
/* ===========================================
日付の確認
ARGS
	date:[string]:検査する日付文字列
RETURN
	検索結果(正:true/誤:false)
=========================================== */
func DateFormatChecker(d string) bool {

	//不要な文字列を削除
	reg := regexp.MustCompile(`[-|/|:| |　]`)
	str := reg.ReplaceAllString(d, "")

	//数値の値に変換するフォーマットを定義
	format := string([]rune("20060102150405")[:len(str)])

	//日付文字列をパースしてエラーが出ないかを確認
	_, err := time.Parse(format, str)
	return err == nil
}

///////////////////////////////////////////////////
/* ===========================================
日付差分
ARGS
	startDate: 開始日
	endDate: 終了日
	delimiter: 日付のデリミタ
RETURN
	日数
=========================================== */
func DateDiffCalculator(startDate, endDate, delimiter string) int {

	sArr := strings.Split(startDate, delimiter)
	sy, _ := strconv.Atoi(sArr[0])
	sm, _ := strconv.Atoi(sArr[1])
	sd, _ := strconv.Atoi(sArr[2])

	eArr := strings.Split(endDate, delimiter)
	ey, _ := strconv.Atoi(eArr[0])
	em, _ := strconv.Atoi(eArr[1])
	ed, _ := strconv.Atoi(eArr[2])

	start := time.Date(sy, time.Month(sm), sd, 0, 0, 0, 0, time.Local)
	end := time.Date(ey, time.Month(em), ed, 0, 0, 0, 0, time.Local)

	//Diff
	diffDays := end.Sub(start).Hours() / 24

	return int(diffDays)
}

///////////////////////////////////////////////////
/* ===========================================
文字列の日付をUTCのTime型に変換
ARGS
	sDateTime: 日時文字列 Foamat:"2023-02-24T11:42:04"
	dtDelimiter: 日付と時間のデリミタ
	dDelimiter: 日付のデリミタ
	tDelimiter: 時間のデリミタ
RETURN
	Time型
=========================================== */
func ConvertStringDateToUtcTime(sDateTime, dtDelimiter, dDelimiter, tDelimiter string, utc bool) time.Time {

	dtArray := strings.Split(sDateTime, dtDelimiter)

	dt := strings.Split(dtArray[0], dDelimiter)
	dy, _ := strconv.Atoi(dt[0])
	dm, _ := strconv.Atoi(dt[1])
	dd, _ := strconv.Atoi(dt[2])

	ts := strings.Split(dtArray[1], tDelimiter)
	tHour, _ := strconv.Atoi(ts[0])
	tMin, _ := strconv.Atoi(ts[1])
	tSec, _ := strconv.Atoi(ts[2])
	tMsec := 0

	var tm time.Time
	if utc {
		tm = time.Date(dy, time.Month(dm), dd, tHour, tMin, tSec, tMsec, time.UTC)
	} else {
		tm = time.Date(dy, time.Month(dm), dd, tHour, tMin, tSec, tMsec, time.Local)
	}

	return tm
}

///////////////////////////////////////////////////
/* ===========================================
日付追加
ARGS
	startDate: 開始日付
	delimiter: 日付のデリミタ
	addD: 追加日数
RETURN
	開始日付から日数分の日付を配列で戻す
=========================================== */
func DateAddCalculator(srcDate, delimiter string, addD int) []string {

	sArr := strings.Split(srcDate, delimiter)
	sy, _ := strconv.Atoi(sArr[0])
	sm, _ := strconv.Atoi(sArr[1])
	sd, _ := strconv.Atoi(sArr[2])

	var dArr []string
	for i := 0; i <= addD; i++ {
		t := time.Date(sy, time.Month(sm), sd, 0, 0, 0, 0, time.Local)
		t = t.AddDate(0, 0, i)
		ta := strings.Split(t.String(), " ")[0]
		ts := strings.Replace(ta, "-", delimiter, -1)
		dArr = append(dArr, ts)
	}

	return dArr
}

///////////////////////////////////////////////////
/* ===========================================
次の日付を取得
ARGS
	frequency	[string]:
		"daily"-->次の日,
		"weekly"-->翌週日曜日の日付,
		"monthly"-->翌月初日,
	startDate	[string]:	起点とする日付文字列
	delimiter	[string]:	日付のデリミタ
	startOrEnd [string]: "s" or "e"
RETURN
	日付	[string]
=========================================== */
func GetNextFrequencyDates(frequency, startDate, delimiter string) []string {
	nextSE := make([]string, 2)
	switch frequency {
	case "daily":
		nextSE[0] = DateAddCalculator(startDate, delimiter, 1)[1] //起点日に1日追加
		nextSE[1] = nextSE[0]
	case "weekly":
		//起点日の曜日を調べて、次の日曜日を検索(日付追加し、最後の配列の要素)
		sArr := strings.Split(startDate, delimiter)
		sy, _ := strconv.Atoi(sArr[0])
		sm, _ := strconv.Atoi(sArr[1])
		sd, _ := strconv.Atoi(sArr[2])
		t := time.Date(sy, time.Month(sm), sd, 0, 0, 0, 0, time.Local)
		nds := DateAddCalculator(startDate, delimiter, 7-int(t.Weekday())) //次の日曜日までの日付
		nextSE[0] = nds[len(nds)-1]
		nde := DateAddCalculator(nextSE[0], delimiter, 6) //次の土曜日までの日付
		nextSE[1] = nde[len(nde)-1]
	case "monthly":
		sArr := strings.Split(startDate, delimiter)
		sy, _ := strconv.Atoi(sArr[0])
		sm, _ := strconv.Atoi(sArr[1])
		var t time.Time
		t = time.Date(sy, time.Month(sm+1), 1, 0, 0, 0, 0, time.Local) //翌月の1日
		nextSE[0] = strings.Split(t.String(), " ")[0]
		t = time.Date(sy, time.Month(sm+2), 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1) //翌々月の初日から1dを引く
		nextSE[1] = strings.Split(t.String(), " ")[0]
	}
	return nextSE
}
