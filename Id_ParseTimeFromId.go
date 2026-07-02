package hgmIdGen

import (
	"time"
)

// 解析普通数据库 id 里的毫秒时间，失败时返回零值时间。
func ParseTimeFromId(id string) time.Time {
	if len(id) != 29 {
		return time.Time{}
	}
	buf, err := _HgmDecodeString(id)
	if err != nil {
		return time.Time{}
	}
	if len(buf) != 18 {
		return time.Time{}
	}
	return parseIdGenMsTime(buf[:6])
}

func parseIdGenMsTime(buf []byte) time.Time {
	if len(buf) != 6 {
		return time.Time{}
	}
	ms := int64(buf[0])<<40 |
		int64(buf[1])<<32 |
		int64(buf[2])<<24 |
		int64(buf[3])<<16 |
		int64(buf[4])<<8 |
		int64(buf[5])
	return time.UnixMilli(ms)
}
