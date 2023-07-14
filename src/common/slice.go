/*
======================
共通処理のファイル
========================
*/
package common

import "regexp"

///////////////////////////////////////////////////
/* ===========================================
スライス割り
ARGS
	sliceLen: スライスの長さ
	devideQty: 割る数
RETURN
	開始位置、終了位置を配列で戻す
=========================================== */
func SliceDevideCalculator(sliceLen, devideQty int) (s, e []int) {
	//v := []int{1,2,3,4}
	//var s []int
	//var e []int

	l := sliceLen
	u := devideQty
	var buffE int
	for i := 0; i < u; i++ {

		if i == 0 {
			s = append(s, 0)
		} else {
			s = append(s, buffE)
		}
		if i == 0 {
			e = append(e, l/u)
		} else if u-1 == i {
			e = append(e, buffE+l-buffE)
		} else {
			e = append(e, l/u+buffE)
		}
		buffE = e[i]
	}

	/* testing
	for i, _ := range s {
		fmt.Println(v[s[i]:e[i]])
	}
	*/
	return s, e
}

///////////////////////////////////////////////////
/* ===========================================
重複を削除する
ARGS
	args: スライス
RETURN
	重複を除外したスライス
=========================================== */
func RemoveDuplicateArrayString(args []string) []string {
	results := make([]string, 0, len(args))
	encountered := map[string]bool{}
	for i := 0; i < len(args); i++ {
		if !encountered[args[i]] {
			encountered[args[i]] = true
			results = append(results, args[i])
		}
	}
	return results
}

///////////////////////////////////////////////////
/* ===========================================
スライス検索
ARGS
	slice: スライス
	value: 検索値
RETURN
	検索結果(あり:true/なし:false)
=========================================== */
func StringSliceSearch(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}
	return false
}

///////////////////////////////////////////////////
/* ===========================================
Stringスライス特定要素の削除
ARGS
	ss: スライス
	ds: 削除したい文字列
RETURN
	特定要素を抜いたスライス
=========================================== */
func RemoveSliceValue(ss []string, ds string) []string {
	result := []string{}
	for _, s := range ss {
		if s == ds {
			continue
		} else {
			result = append(result, s)
		}
	}
	return result
}

///////////////////////////////////////////////////
/* ===========================================
正規表現で一致した要素だけにする
ARGS
	ss: スライス
	re:　正規表現
RETURN
	正規表現に一致した要素のみのスライス
=========================================== */
func MatchedRegExpSliceValue(ss []string, re string) []string {
	result := []string{}

	//正規表現
	reg := regexp.MustCompile(re)
	for _, s := range ss {
		if reg.MatchString(s) {
			result = append(result, s)
		}
	}
	return result
}
