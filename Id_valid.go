package hgmIdGen

// IsValidId 判断一个字符串是否为合法的普通 id（base32编码，长度29）
func IsValidId(id string) bool {
	if len(id) != 29 {
		return false
	}
	return _HgmEncodeIsValid(id)
}