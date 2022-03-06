package camo

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

// テスト用の乱数列モック
type testPatternProvider struct {
	pat []byte
}

func (p *testPatternProvider) Pattern(n byte) []byte {
	return p.pat
}

func (p *testPatternProvider) Length() int {
	return 1
}

// check randomPattern の妥当性チェックを行う
func (p *RandomPattern) check(t *testing.T) {
	m := make(map[string]int)
	for i := 0; i < p.Length(); i++ {
		r := p.Pattern(byte(i))
		// 16進文字列化
		h := hex.EncodeToString(r)
		// 文字列マップで重複の有無をチェック
		_, ok := m[h]
		if ok {
			// 同じ16進文字列があった
			t.Fatalf("%s", h)
		}
		m[h] = len(r)
	}
}

// TestPattern 乱数列パターンのテスト用
func TestPattern(t *testing.T) {
	standardPat().check(t)
	compactPat().check(t)
}

// printRandomPattern 乱数列マトリックスのソースコードを出力
func printRandomPattern(pat RandomPattern) {
	fmt.Println("\treturn &randomPattern{")
	max := pat.Length()
	for i := 0; i < max; i++ {
		nonce := pat.Pattern(byte(i))
		elems := []string{}
		for j := 0; j < len(nonce); j++ {
			elems = append(elems, fmt.Sprintf("%d", nonce[j]))
		}
		fmt.Printf("\t\t{%s},\n", strings.Join(elems, ", "))
	}
	fmt.Println("\t}")
}

// TestMakeRandomPattern4 4bit新規パターンの作成
func TestMakeRandomPattern4(t *testing.T) {
	printRandomPattern(MakeRandomPattern(16/2, []int{53, 59, 61, 67, 71, 73, 79, 83}))
	t.Log()
}

// TestMakeRandomPattern8 8bit新規パターンの作成
func TestMakeRandomPattern8(t *testing.T) {
	printRandomPattern(MakeRandomPattern(256/2, []int{53, 59, 61, 67, 71, 73, 79, 83}))
	t.Log()
}
