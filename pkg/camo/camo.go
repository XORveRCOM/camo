// Package camo 平文迷彩処理を行う
package camo

// PatternProvider 迷彩パターン供給
type PatternProvider interface {
	// Pattern は指定した n 番の乱数列を返します。
	// n は何を指定しても剰余で計算されるためエラーにはなりません。
	// 8ビットの番号であるのは、多ければよいというわけではないからです。
	// 例えば4ビット迷彩処理では0～15の指定のみ使用されます。
	Pattern(n byte) []byte
	// Length は迷彩パターンのパターン数を返します。
	Length() int
}

// MixingContext 迷彩処理コンテキスト
// 平文を迷彩設定して迷彩平文にする Mixing() と、それを戻す UnMixing() で構成されます。
// 一度にすべてのバイト列を渡さなくとも、逐次に変換処理ができます。
// スレッド安全ではありませんので、複数のスレッドで使い回さないようにしてください。
type MixingContext interface {
	// Mixing 迷彩設定
	Mixing(pt []byte) []byte
	// UnMixing 迷彩解除
	UnMixing(pt []byte) []byte
}
