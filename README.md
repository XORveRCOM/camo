# camo

## 特許情報

日本ではJP-5992651として取得済みである特許の参考実装である。
<https://patents.google.com/patent/JP5992651B2/ja>

本ライブラリはあくまでプロプライエタリなソフトウェアであり、現時点では評価のためだけに供されるものとする。  

## 実行例

<https://play.golang.org/p/aX9NsCHkBn7>

```golang
package main

import (
    "fmt"
    "github.com/xorvercom/camo/pkg/camo"
)

func main() {
    b := []byte("abcd0123")
    fmt.Printf("plane: %v\n", b)
    // 迷彩設定
    ctxEnc := camo.NewStandardCamo()
    cb := ctxEnc.Mixing(b)
    fmt.Printf("camo: %v\n", cb)
    // 迷彩解除
    ctxDec := camo.NewStandardCamo()
    fmt.Println(string(ctxDec.UnMixing(cb)))
}
```

<https://play.golang.org/p/VpMnuAgOWlU>

```golang
package main

import (
    "fmt"
    "github.com/xorvercom/camo/pkg/camo"
)

func main() {
    b := []byte("abcd0123")
    fmt.Printf("plane: %v\n", b)
    // 迷彩設定
    ctxEnc := camo.NewCompactCamo()
    cb := ctxEnc.Mixing(b)
    fmt.Printf("camo: %v\n", cb)
    // 迷彩解除
    ctxDec := camo.NewCompactCamo()
    fmt.Println(string(ctxDec.UnMixing(cb)))
}
```

&copy;XORveR.com
