// Package camo 平文迷彩処理を行います
package camo

// compactContext 4ビット迷彩処理の管理構造体
type compactContext struct {
	// 迷彩位置
	pt int
	// 迷彩最終位置
	ptend int
	// 迷彩パターン
	pat []byte
	// 迷彩パターン長
	patlen int
	// 乱数列供給
	patProv PatternProvider
	// 乱数生成器
	nonce NonceProvider
}

// NewCompactCamo 4ビット迷彩処理コンテキストの新規作成
// 4ビット迷彩処理では16種類の迷彩平文を生成します。
// 例えば128ビット暗号で32個までの選択平文が使用可能範囲内となってしまいますが、知られている線形解読法では数千個オーダーの選択平文を要するため問題はありません。
func NewCompactCamo() MixingContext {
	return &compactContext{patProv: NewCompactPatternProvider(), nonce: NewNonce()}
}

// setup 4ビット迷彩処理の管理構造体の初期化
func (m *compactContext) setup(val byte) {
	n := val >> (8 - 4)
	l := val & 0x0f
	m.pat = m.patProv.Pattern(n)
	m.patlen = len(m.pat)
	m.pt = 0
	m.ptend = int(l) * m.patlen
}

// Mixing 4ビット迷彩設定処理
func (m *compactContext) Mixing(pt []byte) []byte {
	if m == nil {
		// nilレシーバはコーディングミス
		panic("nil.Mixing()")
	}
	ret := []byte{}
	if pt == nil || len(pt) == 0 {
		return ret
	}
	for i := 0; i < len(pt); i++ {
		for m.pt == m.ptend {
			// リセット
			*m = compactContext{patProv: m.patProv, nonce: m.nonce}
			n := uint8(m.nonce.Nonce())
			m.setup(n)
			// ヘッダの埋め込み
			ret = append(ret, n)
		}
		// マスク作成
		mask := m.pat[m.pt%m.patlen]
		m.pt++
		// XOR
		ret = append(ret, pt[i]^mask)
	}
	return ret
}

// UnMixing 4ビット迷彩解除処理
func (m *compactContext) UnMixing(pt []byte) []byte {
	if m == nil {
		// nilレシーバはコーディングミス
		panic("nil.UnMixing()")
	}
	ret := []byte{}
	if pt == nil || len(pt) == 0 {
		return ret
	}
	for _, b := range pt {
		if m.pt == m.ptend {
			// 切替
			m.setup(b)
		} else {
			// マスク作成
			mask := m.pat[m.pt%m.patlen]
			m.pt++
			// XOR
			ret = append(ret, b^mask)
		}
	}
	return ret
}
