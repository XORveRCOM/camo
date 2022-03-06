package camo

type RandomPattern [][]byte

// MakeRandomPattern 乱数列マトリックスを新規作成
func MakeRandomPatternP(halfCnt int, s []int) RandomPattern {
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

// // MakeRandomPattern 乱数列マトリックスを新規作成
// func MakeRandomPattern(halfCnt int, s []int) RandomPattern {
// 	var nonces = make([]byte, 1)
// 	// 乱数列の生成
// 	ret := [][]byte{}
// 	for i := 0; i < halfCnt; i++ {
// 		ran := []byte{}
// 		inv := []byte{}
// 		for j := 0; j < s[i%len(s)]; j++ {
// 			// 安全な乱数を使う
// 			cnt, err := crand.Read(nonces)
// 			if cnt != cap(nonces) {
// 				panic(fmt.Errorf("rand.Read() : cnt=%d", cnt))
// 			} else if err != nil {
// 				panic(fmt.Errorf("rand.Read() : %v", err))
// 			}
// 			val := uint8(nonces[0])
// 			ran = append(ran, val)
// 			inv = append(inv, ^val)
// 		}
// 		ret = append(ret, ran, inv)
// 	}
// 	// シャッフル
// 	random := NewNonce()
// 	max := len(ret)
// 	for i := 0; i < halfCnt*2; i++ {
// 		from := int(random.Nonce()) % max
// 		to := int(random.Nonce()) % max
// 		ret[to], ret[from] = ret[from], ret[to]
// 	}
// 	return ret
// }

// pattern 乱数列をコピーして返す
func (p *RandomPattern) Pattern(n byte) []byte {
	//pat := (*p)[n]
	//ret := make([]byte, len(pat))
	//_ = copy(ret, pat)
	//return ret
	return (*p)[int(n)%len(*p)]
}

// Length 乱数列マトリックスのインデックス数
func (p *RandomPattern) Length() int {
	return len(*p)
}
