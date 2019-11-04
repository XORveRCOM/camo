package camo

import (
	"encoding/hex"
	"fmt"
	"testing" // テストユーティリティパッケージをimport
)

func (m *standardContext) String() string {
	return fmt.Sprintf("{stage:%d, n:%02x, l:%d, pt:%d, ptend:%d, patlen:%d}", m.stage, m.n, m.l, m.pt, m.ptend, m.patlen)
}

// TestStandardMixing 迷彩設定テスト
func TestStandardMixing(t *testing.T) {
	var mask byte = 255

	// コンテキスト生成
	contextEnc := NewStandardCamo()
	ctx, ok := contextEnc.(*standardContext)
	if !ok {
		t.Error("NewStandardCamo() error")
	}
	// テスト用乱数列に差し替え
	ctx.patProv = &testPatternProvider{pat: []byte{mask}}
	// テスト用乱数生成器に差し替え
	rn := []int64{0x0000, 0x0101, 0x0202, 0x0303, 0x0404, 0x0505}
	ctx.nonce = &testNonceProvider{buf: rn}
	t.Logf("%v", ctx)

	// 平文
	src := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	srcstr := hex.EncodeToString(src)

	// 迷彩平文
	mix := ctx.Mixing(src)
	// 手計算で作る迷彩平文
	check := []byte{}
	check = append(check, byte(rn[0]>>0), byte(rn[0]>>8))
	// l:0 の場合は即座に次のヘッダ
	check = append(check, byte(rn[1]>>0), byte(rn[1]>>8))
	check = append(check, src[0]^mask)
	check = append(check, byte(rn[2]>>0), byte(rn[2]>>8))
	check = append(check, src[1]^mask, src[2]^mask)
	check = append(check, byte(rn[3]>>0), byte(rn[3]>>8))
	check = append(check, src[3]^mask, src[4]^mask, src[5]^mask)
	check = append(check, byte(rn[4]>>0), byte(rn[4]>>8))
	check = append(check, src[6]^mask, src[7]^mask, src[8]^mask, src[9]^mask)
	check = append(check, byte(rn[5]>>0), byte(rn[5]>>8))
	check = append(check, src[10]^255)

	// 比較
	mixstr := hex.EncodeToString(mix)
	checkstr := hex.EncodeToString(check)
	if checkstr != mixstr {
		t.Fatalf("Mixing fail \n\tmix:%s\n\tchk:%s", mixstr, checkstr)
	}

	// 迷彩解除して比較
	contextDec := NewStandardCamo()
	ctxDec, ok := contextDec.(*standardContext)
	if !ok {
		t.Error("NewStandardCamo() error")
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

// TestStandardUnMixing 迷彩解除テスト
func TestStandardUnMixing(t *testing.T) {
	// テスト用バイト列
	src := make([]byte, 1024*1024)
	srcHex := hex.EncodeToString(src)

	// 迷彩設定
	contextEnc := NewStandardCamo()
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
	contextDec := NewStandardCamo()
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

// TestStandardMixingError 迷彩処理のエラー
func TestStandardMixingError(t *testing.T) {
	ctx := NewStandardCamo()
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
		var ctx2 *standardContext
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
		var ctx2 *standardContext
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

// // ラッシュテスト
// func TestMM512(t *testing.T) {
// 	for i := 0; i < 512; i++ {
// 		TestUnMixing(t)
// 	}
// }
