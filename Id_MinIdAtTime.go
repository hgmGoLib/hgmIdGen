package hgmIdGen

import (
	"time"
)

// 生成指定毫秒时刻对应的最小普通数据库 id，适合做数据库时间范围扫描下界。
// 最小 id 即: 6 字节时间 + 后续递增和随机段全 0, 再 base32 编码.
func MinIdAtTime(t time.Time) string {
	ms := uint64(t.UnixMilli())
	var dst [29]byte
	var src [18]byte
	src[0] = byte(ms >> 40)
	src[1] = byte(ms >> 32)
	src[2] = byte(ms >> 24)
	src[3] = byte(ms >> 16)
	src[4] = byte(ms >> 8)
	src[5] = byte(ms)
	_HgmEncodeSlice(dst[:], src[:])
	return string(dst[:])
}
