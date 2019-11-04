package camo

// テスト用の乱数生成器モック
type testNonceProvider struct {
	pt  int
	buf []int64
}

// テスト用のモック testNonceProvider が NonceProvider に適応するようにする
func (r *testNonceProvider) Nonce() int64 {
	b := r.buf[r.pt%len(r.buf)]
	r.pt++
	return b
}
