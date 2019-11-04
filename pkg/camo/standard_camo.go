// Package camo 平文迷彩処理を行います
package camo

// standardContext 迷彩処理の管理構造体
type standardContext struct {
	// 初期化段階
	stage int
	// NN
	n byte
	// 迷彩繰り返し回数
	l byte
	// 迷彩位置
	pt int
	// 迷彩最終位置+1
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

// NewStandardCamo 迷彩処理コンテキストの新規作成
func NewStandardCamo() MixingContext {
	return &standardContext{patProv: NewStandardPatternProvider(), nonce: NewNonce()}
}

// setup 迷彩処理の管理構造体の初期化
func (m *standardContext) setup(val byte) {
	if m.stage == 0 {
		m.stage = 1
		m.n = val
		m.pat = m.patProv.Pattern(m.n)
		m.patlen = len(m.pat)
		m.pt = 0
		m.ptend = 0
	} else if m.stage == 1 {
		m.stage = 0
		m.l = val
		m.ptend = (int(val)) * m.patlen
	}
}

// Mixing 迷彩設定処理
func (m *standardContext) Mixing(pt []byte) []byte {
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
			*m = standardContext{patProv: m.patProv, nonce: m.nonce}
			n := m.nonce.Nonce()
			m.setup(uint8(n >> 0))
			m.setup(uint8(n >> 8))
			// ヘッダの埋め込み
			ret = append(ret, m.n, m.l)
		}
		// マスク作成
		mask := m.pat[m.pt%m.patlen]
		m.pt++
		// XOR
		ret = append(ret, pt[i]^mask)
	}
	return ret
}

// UnMixing 迷彩解除処理
func (m *standardContext) UnMixing(pt []byte) []byte {
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
