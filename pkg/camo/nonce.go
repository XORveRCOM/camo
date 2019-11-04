package camo

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
	"sync"
)

var (
	// createRandom() 用MUTEX
	createRandomMutex *sync.Mutex = new(sync.Mutex)
)

// NewNonce 乱数値生成器を作成
func NewNonce() NonceProvider {
	return &nonce{random: createRandom()}
}

// NonceProvider 乱数値を返すインタフェース
type NonceProvider interface {
	// Nonce は int64 の疑似乱数値を返します。
	Nonce() int64
}

type nonce struct {
	random *mrand.Rand
}

// Nonce ナンス生成
func (r *nonce) Nonce() int64 {
	return r.random.Int63()
}

// 乱数生成器
func createRandom() *mrand.Rand {
	// crypto/rand はスレッドセーフではない疑いがある (Go 1.13)
	// その理由は例えば Windows では CryptGenRandom() に渡す HCRYPTPROV をプロセス内で使いまわしている
	// CSP のモジュールのハンドルということで大丈夫だとは思われる
	// しかし、そんな心配をするくらいならば、自前で排他してしまう方が安心
	createRandomMutex.Lock()
	defer createRandomMutex.Unlock()

	var iseed int64
	var seed = make([]byte, 8)
	// シードは安全な乱数を使う
	cnt, err := crand.Read(seed)
	if cnt != cap(seed) {
		panic(fmt.Errorf("rand.Read() : cnt=%d", cnt))
	} else if err != nil {
		panic(fmt.Errorf("rand.Read() : %v", err))
	}
	// シードから疑似乱数を生成
	buf := bytes.NewReader(seed)
	binary.Read(buf, binary.LittleEndian, &iseed)
	return mrand.New(mrand.NewSource(iseed))
}
