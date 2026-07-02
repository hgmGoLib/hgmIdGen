package hgmIdGen

// 本文件从 hgmLib/hgmEncode/hgmEncodeB32 复制而来, 使本库独立, 不依赖 hgmLib.

import (
	"encoding/base32"
	"sync"
)

func _HgmDecodeString(s string) ([]byte, error) {
	_hgmInit()
	return gHgmEncode.DecodeString(s)
}

func _HgmEncodeSlice(dst []byte, src []byte) {
	_hgmInit()
	gHgmEncode.Encode(dst, src)
}

func _hgmInit() {
	gHgmEncodeOnce.Do(func() {
		gHgmEncode = base32.NewEncoding(hgmEncodingS).WithPadding(base32.NoPadding)
	})
}

const hgmEncodingS = "123456789abcdefghjkmnpqrstuvwxyz"

var gHgmEncode *base32.Encoding
var gHgmEncodeOnce sync.Once

var validCharSet map[byte]struct{}
var validCharSetOnce sync.Once

func _initValidCharSet() {
	validCharSetOnce.Do(func() {
		validCharSet = make(map[byte]struct{}, 32)
		for i := 0; i < len(hgmEncodingS); i++ {
			validCharSet[hgmEncodingS[i]] = struct{}{}
		}
	})
}

func _HgmEncodeIsValid(s string) bool {
	_initValidCharSet()
	for i := 0; i < len(s); i++ {
		_, ok := validCharSet[s[i]]
		if !ok {
			return false
		}
	}
	return true
}
