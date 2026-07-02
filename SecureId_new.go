package hgmIdGen

import (
	"time"
)

/*
输出结果例子 17fatweptd1swfgg11ekgtx5j459wnf112v7jp81wn
字符串固定长度 42
高安全 id
==========
* 适合高安全性要求的id生成. 比如 sessionId 和 acceessToken 有人故意ddos碰撞也很难碰撞上.
* 固定 26个字节二进制.
* 6个字节 时间到毫秒. ( 大约 8925 年)
* 1-4个字节 递增变长编码 当前进程挂锁顺序增加id (时间变了之后,从0 开始递增) (2位表示后面有几个字节,6-30位表示后面的id 最大 2^30-1 = 1073741823)
* 16-19个字节 随机二进制数据.(来自强随机api)
* 然后base32 编码 使用不易混淆的 字符,使用按ascii码递增排列的字符串
 */
func NewSecureId() string {
	ms := uint64(time.Now().UnixMilli())
	incrU32 := nextId(ms)
	return secureEncode(ms, incrU32)
}

func secureEncode(ms uint64, incrU32 uint32) string {
	var dst [42]byte 
	var src [26]byte
	// 时间戳 6字节
	src[0] = byte(ms >> 40)
	src[1] = byte(ms >> 32)
	src[2] = byte(ms >> 24)
	src[3] = byte(ms >> 16)
	src[4] = byte(ms >> 8)
	src[5] = byte(ms)
	// 递增 ID 编码（变长）
	incrLen := encodeVarLenUint32(src[6:], incrU32)
	// 随机数据（16~19字节）
	randOffset := 6 + incrLen
	gRf.MustReadWithLock(src[randOffset : ])
	// base32 编码（使用自定义安全字符表）
	_HgmEncodeSlice(dst[:], src[:])
	return string(dst[:])
}