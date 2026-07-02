package hgmIdGen

import (
	"errors"
	"strconv"
	"time"
)

// 生成18字节固定长度的二进制id,方案和NewId()完全一致,只是跳过了base32编码步骤.
// 传入dst用于复用内存分配,结果append到dst后面返回.
func NewIdBinary(dst []byte) []byte {
	ms := uint64(time.Now().UnixMilli())
	incrU32 := nextId(ms)
	start := len(dst)
	dst = append(dst, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	fillIdBytes(dst[start:start+18], ms, incrU32)
	return dst
}

// 把18字节二进制id转换为29字符的base32字符串id,和NewId()的输出格式一致.
func IdBinaryToId(bin []byte) (string, error) {
	if len(bin) != 18 {
		return "", errors.New("hgmIdGen.IdBinaryToId: invalid bin length " + strconv.Itoa(len(bin)) + ", expect 18")
	}
	var dst [29]byte
	_HgmEncodeSlice(dst[:], bin)
	return string(dst[:]), nil
}

// 把29字符的base32字符串id转换为18字节二进制id.
func IdToIdBinary(id string) ([]byte, error) {
	if len(id) != 29 {
		return nil, errors.New("hgmIdGen.IdToIdBinary: invalid id length " + strconv.Itoa(len(id)) + ", expect 29")
	}
	if !_HgmEncodeIsValid(id) {
		return nil, errors.New("hgmIdGen.IdToIdBinary: invalid base32 characters in id")
	}
	bin, err := _HgmDecodeString(id)
	if err != nil {
		return nil, errors.New("hgmIdGen.IdToIdBinary: base32 decode failed: " + err.Error())
	}
	if len(bin) != 18 {
		return nil, errors.New("hgmIdGen.IdToIdBinary: decoded length " + strconv.Itoa(len(bin)) + ", expect 18")
	}
	return bin, nil
}
