package camo

import (
	crand "crypto/rand"
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
func (p *randomPattern) check(t *testing.T) {
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

// makeRandomPattern 乱数列マトリックスを新規作成
func makeRandomPatternP(halfCnt int, s []int) randomPattern {
	// 乱数列の生成
	random := NewNonce()
	ret := [][]byte{}
	for i := 0; i < halfCnt; i++ {
		ran := []byte{}
		inv := []byte{}
		for j := 0; j < s[i%len(s)]; j++ {
			val := uint8(random.Nonce())
			ran = append(ran, val)
			inv = append(inv, ^val)
		}
		ret = append(ret, ran, inv)
	}
	// シャッフル
	max := len(ret)
	for i := 0; i < halfCnt*2; i++ {
		from := int(random.Nonce()) % max
		to := int(random.Nonce()) % max
		ret[to], ret[from] = ret[from], ret[to]
	}
	return ret
}

// makeRandomPattern 乱数列マトリックスを新規作成
func makeRandomPattern(halfCnt int, s []int) randomPattern {
	var nonces = make([]byte, 1)
	// 乱数列の生成
	ret := [][]byte{}
	for i := 0; i < halfCnt; i++ {
		ran := []byte{}
		inv := []byte{}
		for j := 0; j < s[i%len(s)]; j++ {
			// 安全な乱数を使う
			cnt, err := crand.Read(nonces)
			if cnt != cap(nonces) {
				panic(fmt.Errorf("rand.Read() : cnt=%d", cnt))
			} else if err != nil {
				panic(fmt.Errorf("rand.Read() : %v", err))
			}
			val := uint8(nonces[0])
			ran = append(ran, val)
			inv = append(inv, ^val)
		}
		ret = append(ret, ran, inv)
	}
	// シャッフル
	random := NewNonce()
	max := len(ret)
	for i := 0; i < halfCnt*2; i++ {
		from := int(random.Nonce()) % max
		to := int(random.Nonce()) % max
		ret[to], ret[from] = ret[from], ret[to]
	}
	return ret
}

// printRandomPattern 乱数列マトリックスのソースコードを出力
func printRandomPattern(pat randomPattern) {
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
	printRandomPattern(makeRandomPattern(16/2, []int{53, 59, 61, 67, 71, 73, 79, 83}))
	t.Log()
}

// TestMakeRandomPattern8 8bit新規パターンの作成
func TestMakeRandomPattern8(t *testing.T) {
	printRandomPattern(makeRandomPattern(256/2, []int{53, 59, 61, 67, 71, 73, 79, 83}))
	t.Log()
}
