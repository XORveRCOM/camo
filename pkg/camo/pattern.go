package camo

type randomPattern [][]byte

// pattern 乱数列をコピーして返す
func (p *randomPattern) Pattern(n byte) []byte {
	//pat := (*p)[n]
	//ret := make([]byte, len(pat))
	//_ = copy(ret, pat)
	//return ret
	return (*p)[int(n)%len(*p)]
}

// Length 乱数列マトリックスのインデックス数
func (p *randomPattern) Length() int {
	return len(*p)
}
