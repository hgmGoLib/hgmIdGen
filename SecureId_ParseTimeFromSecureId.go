package hgmIdGen

import (
	"time"
)

// 解析高安全 id 里的毫秒时间，失败时返回零值时间。
func ParseTimeFromSecureId(id string) time.Time {
	if len(id) != 42 {
		return time.Time{}
	}
	buf, err := _HgmDecodeString(id)
	if err != nil {
		return time.Time{}
	}
	if len(buf) != 26 {
		return time.Time{}
	}
	return parseIdGenMsTime(buf[:6])
}
