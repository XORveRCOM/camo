package camo

import (
	"encoding/hex"
	"fmt"
	"testing" // テストユーティリティパッケージをimport
)

func (m *compactContext) String() string {
	return fmt.Sprintf("{pt:%d, ptend:%d, patlen:%d}", m.pt, m.ptend, m.patlen)
}

// TestCompactMixing 迷彩設定テスト
func TestCompactMixing(t *testing.T) {
	var mask byte = 255

	// コンテキスト生成
	contextEnc := NewCompactCamo()
	ctx, ok := contextEnc.(*compactContext)
	if !ok {
		t.Error("NewCompactCamo() error")
	}
	// テスト用乱数列に差し替え
	ctx.patProv = &testPatternProvider{pat: []byte{mask}}
	// テスト用乱数生成器に差し替え
	rn := []int64{(0 << 4) + 0, (1 << 4) + 1, (2 << 4) + 2, (3 << 4) + 3, (4 << 4) + 4, (5 << 4) + 5}
	ctx.nonce = &testNonceProvider{buf: rn}
	t.Logf("%v", ctx)

	// 平文
	src := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	srcstr := hex.EncodeToString(src)

	// 迷彩平文
	mix := ctx.Mixing(src)
	// 手計算で作る迷彩平文
	check := []byte{}
	check = append(check, byte(rn[0]))
	// l:0 の場合は即座に次のヘッダ
	check = append(check, byte(rn[1]))
	check = append(check, src[0]^mask)
	check = append(check, byte(rn[2]))
	check = append(check, src[1]^mask, src[2]^mask)
	check = append(check, byte(rn[3]))
	check = append(check, src[3]^mask, src[4]^mask, src[5]^mask)
	check = append(check, byte(rn[4]))
	check = append(check, src[6]^mask, src[7]^mask, src[8]^mask, src[9]^mask)
	check = append(check, byte(rn[5]))
	check = append(check, src[10]^mask)

	// 比較
	mixstr := hex.EncodeToString(mix)
	checkstr := hex.EncodeToString(check)
	if checkstr != mixstr {
		t.Fatalf("Mixing fail \n\tmix:%s\n\tchk:%s", mixstr, checkstr)
	}

	// 迷彩解除して比較
	contextDec := NewCompactCamo()
	ctxDec, ok := contextDec.(*compactContext)
	if !ok {
		t.Error("NewCompactCamo() error")
	}
	// テスト用乱数列に差し替え
	ctxDec.patProv = &testPatternProvider{pat: []byte{mask}}
	// テスト用乱数生成器に差し替え
	ctxDec.nonce = &testNonceProvider{buf: rn}
	t.Logf("%v", ctxDec)

	dst := ctxDec.UnMixing(mix)
	dststr := hex.EncodeToString(dst)
	if srcstr != dststr {
		t.Fatalf("UnMixing fail \n\tdst:%s\n\tsrc:%s", dststr, srcstr)
	}
}

// TestCompactUnMixing 迷彩解除テスト
func TestCompactUnMixing(t *testing.T) {
	// テスト用バイト列
	src := make([]byte, 1024*1024)
	srcHex := hex.EncodeToString(src)

	// 迷彩設定
	contextEnc := NewCompactCamo()
	ptm := []byte{}
	for i := 0; i < len(src); i++ {
		// バイト単位で小分けに迷彩設定（ストリームとして操作するテスト）
		ptm = append(ptm, contextEnc.Mixing(src[i:i+1])...)
	}
	ptmHex := hex.EncodeToString(ptm)
	if srcHex == ptmHex {
		t.Fatalf("src == ptm: %s\n", srcHex)
	}

	// 迷彩解除
	contextDec := NewCompactCamo()
	pt := []byte{}
	for i := 0; i < len(ptm); i++ {
		// バイト単位で小分けに迷彩解除（ストリームとして操作するテスト）
		pt = append(pt, contextDec.UnMixing(ptm[i:i+1])...)
	}
	dstHex := hex.EncodeToString(pt)
	if srcHex != dstHex {
		t.Fatalf("src:%s\ndst: %s\n", srcHex, dstHex)
	}

	t.Logf("src:%d head:%d\n", len(ptm), len(ptm)-len(src))
}

// TestCompactMixingError 迷彩処理のエラー
func TestCompactMixingError(t *testing.T) {
	ctx := NewCompactCamo()
	var s []byte
	s = ctx.Mixing(nil)
	if s == nil || 0 != len(s) {
		t.Fatalf("ctx.Mixing(nil) return %v", s)
	}
	s = ctx.UnMixing(nil)
	if s == nil || 0 != len(s) {
		t.Fatalf("ctx.UnMixing(nil) return %v", s)
	}

	func() {
		var ctx2 *compactContext
		ctx2 = nil
		defer func() {
			err := recover()
			if err != "nil.Mixing()" {
				t.Fatalf("got %v\nwant %v", err, "nil.Mixing()")
			}
		}()
		ctx2.Mixing([]byte{})
	}()

	func() {
		var ctx2 *compactContext
		ctx2 = nil
		defer func() {
			err := recover()
			if err != "nil.UnMixing()" {
				t.Fatalf("got %v\nwant %v", err, "nil.UnMixing()")
			}
		}()
		ctx2.UnMixing([]byte{})
	}()
}
