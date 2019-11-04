package camo

import (
	"fmt"
)

// Example サンプル
func Example() {
	var ctx MixingContext
	var c []byte
	var p []byte

	// 平文を迷彩平文に変換します
	ctx = NewStandardCamo()
	c = []byte{}
	c = append(c, ctx.Mixing([]byte("abcd"))...)
	c = append(c, ctx.Mixing([]byte("0123"))...)

	// 迷彩平文を平文に戻します
	ctx = NewStandardCamo()
	p = ctx.UnMixing(c)
	fmt.Println(string(p))

	// 平文を迷彩平文に変換します
	ctx = NewCompactCamo()
	c = []byte{}
	c = append(c, ctx.Mixing([]byte("abcd"))...)
	c = append(c, ctx.Mixing([]byte("0123"))...)

	// 迷彩平文を平文に戻します
	ctx = NewCompactCamo()
	p = ctx.UnMixing(c)
	fmt.Println(string(p))

	// Output:
	// abcd0123
	// abcd0123
}
