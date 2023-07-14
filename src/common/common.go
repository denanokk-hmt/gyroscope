/*
======================
共通処理のファイル
========================
*/
package common

import (
	"math/rand"
	"strconv"

	"github.com/google/uuid"
)

///////////////////////////////////////////////////
/* ===========================================
比較演算
ARGS
	ls: leftSideの値
	rs: rightSideの値
	ope: オペレーター
RETURN
	条件に一致:true/条件に不一致:false
=========================================== */
func Comparison(ls, rs, ope string) bool {

	var judge bool
	li, _ := strconv.Atoi(ls)
	ri, _ := strconv.Atoi(rs)

	//eq or ne以外は、Filter値を数値に変換した値で比較する
	switch ope {
	case "eq":
		if ls == rs {
			judge = true
		} else {
			judge = false
		}
	case "ne":
		if ls != rs {
			judge = true
		} else {
			judge = false
		}
	case "lt":
		if li < ri {
			judge = true
		} else {
			judge = false
		}
	case "gt":
		if li > ri {
			judge = true
		} else {
			judge = false
		}
	case "le":
		if li <= ri {
			judge = true
		} else {
			judge = false
		}
	case "ge":
		if li >= ri {
			judge = true
		} else {
			judge = false
		}
	}
	return judge
}

///////////////////////////////////////////////////
/* ===========================================
Chunkの開始と終了を計算
ARGS
	length: Chunkさせる配列の要素数
	chunkQty: chunkさせる数
	int: Chunk処理中のLoop index
	start: 配列から抽出する開始位置
	end:　配列から抽出する終了位置
	connter: Chunk数のカウンタ
RETURN
	Chunkする:true/Chunkしない:false
=========================================== */
func ChunkCalculator(length, chunkQty, idx int, start, end, counter *int) bool {

	var judge bool

	//Chunk数450でEntityを区切ってDSを更新していく(=DSの最大Update数500)
	if *counter == chunkQty {
		*counter = 0
	}
	if *counter == 0 {
		*start = idx
		*end = idx + chunkQty
		if *end > length {
			*end = length
		}
		judge = true
	}
	return judge
}

///////////////////////////////////////////////////
/* ===========================================
Chunkの開始と終了を計算
=========================================== */
type Chunks struct {
	Positions []Chunk
}
type Chunk struct {
	Start int
	End   int
	Qty   int
}

///////////////////////////////////////////////////
/* ===========================================
Chunkの開始と終了を計算
ARGS
	length: Chunkさせる配列の要素数
	chunkQty: chunkさせる数
RETURN
	Chunkする配列の抽出開始と終了位置
=========================================== */
func ChunkCalculator2(eLen, chunkQty int) Chunks {

	var counter, start, end int
	var chunks Chunks
	for idx := 0; idx < eLen; idx++ {
		if counter == chunkQty {
			counter = 0
		}
		if counter == 0 {
			start = idx
			end = idx + chunkQty
			if end > eLen {
				end = eLen
			}
			var chunk Chunk = Chunk{
				Start: start,
				End:   end,
				Qty:   end - start}
			chunks.Positions = append(chunks.Positions, chunk)
		}
		counter += 1
	}

	return chunks
}

///////////////////////////////////////////////////
/* ===========================================
UUIv4Dを生成
ARGS
RETURN
	uuidv4 string
	err nil or error
=========================================== */
func CreateUuidV4() (string, error) {
	uuidv4, err := uuid.NewRandom()
	return uuidv4.String(), err
}

///////////////////////////////////////////////////
/* ===========================================
ランダム文字列を生成
ARGS
	num [int] 出力文字数
RETURN
	 string
=========================================== */
func RandomString(n int) string {
	var letter = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
