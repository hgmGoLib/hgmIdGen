package hgmIdGen

// IsValidSecureId 判断一个字符串是否为合法的高安全 id（base32编码，长度42）
func IsValidSecureId(id string) bool {
	if len(id) != 42 {
		return false
	}
	return _HgmEncodeIsValid(id)
}
