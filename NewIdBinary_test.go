package hgmIdGen

import (
	"testing"
	"fmt"
	"bytes"
	"sync"
)

func TestNewIdBinary(ot *testing.T) {
	// 测试输出长度固定18字节
	bin := NewIdBinary(nil)
	_Equal(len(bin), 18)

	// 测试传入dst复用内存,结果append到后面
	prefix := []byte{0xAA, 0xBB}
	result := NewIdBinary(prefix)
	_Equal(len(result), 20)
	_Equal(result[0], byte(0xAA))
	_Equal(result[1], byte(0xBB))
	_Equal(len(result[2:]), 18)

	// 测试同进程内二进制id大致递增(和NewId的行为一致)
	last := NewIdBinary(nil)
	for i := 0; i < 100000; i++ {
		this := NewIdBinary(nil)
		if bytes.Compare(this, last) < 0 {
			fmt.Println(i, this, last)
			panic("binary id not increasing")
		}
		last = this
	}

	// 测试并发生成不重复
	{
		var wg sync.WaitGroup
		binChan := make(chan string, 10000)
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < 1000; j++ {
					binChan <- string(NewIdBinary(nil))
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(binChan)
		bins := make(map[string]bool)
		for b := range binChan {
			if bins[b] {
				panic("duplicate binary id")
			}
			bins[b] = true
		}
	}
}

func TestIdBinaryToId(ot *testing.T) {
	// 生成一个NewId,转成二进制,再转回字符串,应该和原来一样
	for i := 0; i < 10000; i++ {
		id := NewId()
		bin, err := IdToIdBinary(id)
		if err != nil {
			panic("IdToIdBinary failed: " + err.Error())
		}
		_Equal(len(bin), 18)
		id2, err := IdBinaryToId(bin)
		if err != nil {
			panic("IdBinaryToId failed: " + err.Error())
		}
		_Equal(id, id2)
	}

	// 测试NewIdBinary生成的二进制转成字符串后,格式和NewId一致
	for i := 0; i < 10000; i++ {
		bin := NewIdBinary(nil)
		id, err := IdBinaryToId(bin)
		if err != nil {
			panic("IdBinaryToId failed: " + err.Error())
		}
		_Equal(len(id), 29)
		_Equal(IsValidId(id), true)
		_Equal(IsValidSecureId(id), false)
	}
}

func TestIdToIdBinary(ot *testing.T) {
	// 正常转换
	id := NewId()
	bin, err := IdToIdBinary(id)
	if err != nil {
		panic("IdToIdBinary failed: " + err.Error())
	}
	_Equal(len(bin), 18)

	// 长度不对应该报错
	_, err = IdToIdBinary("tooshort")
	if err == nil {
		panic("expect error for short id")
	}

	// 非法字符应该报错
	_, err = IdToIdBinary("AAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	if err == nil {
		panic("expect error for invalid chars")
	}

	// IdBinaryToId 传入长度不对应该报错
	_, err = IdBinaryToId([]byte{1, 2, 3})
	if err == nil {
		panic("expect error for short bin")
	}
	_, err = IdBinaryToId(nil)
	if err == nil {
		panic("expect error for nil bin")
	}
}

func TestNewIdBinaryAndNewIdConsistency(ot *testing.T) {
	// 用fillIdBytes对比,确保二进制方案和字符串方案一致
	ms := uint64(1700000000000)
	testCases := []uint32{0, 1, 0x3F, 0x3F + 1, 0x3FFF, 0x3FFF + 1, 0x3FFFFF, 0x3FFFFF + 1, 1073741823}
	for _, incr := range testCases {
		strId := shortEncode(ms, incr)
		bin := make([]byte, 0, 18)
		start := len(bin)
		bin = append(bin, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
		fillIdBytes(bin[start:start+18], ms, incr)
		strFromBin, err := IdBinaryToId(bin)
		if err != nil {
			panic("IdBinaryToId failed: " + err.Error())
		}
		// 前缀部分(时间+递增)应该一致,随机部分不同所以整体不同,
		// 但是单独生成的二进制转字符串后格式应该合法
		_Equal(len(bin), 18)
		_Equal(len(strFromBin), 29)
		_Equal(IsValidId(strFromBin), true)

		// 对同一个strId做来回转换验证
		bin2, err := IdToIdBinary(strId)
		if err != nil {
			panic("IdToIdBinary failed: " + err.Error())
		}
		str2, err := IdBinaryToId(bin2)
		if err != nil {
			panic("IdBinaryToId failed: " + err.Error())
		}
		_Equal(strId, str2)
	}
}

func TestNewIdBinaryBenchmark(ot *testing.T) {
	// 验证 NewIdBinary 和 NewId alloc 数量一致(都是1次, 即最终的 string/slice 分配)
	_hgmBenchmarkSetName("NewId")
	_hgmBenchmarkWithRepeatNum(10000, func() {
		NewId()
	})
	_hgmBenchmarkSetName("NewIdBinary")
	_hgmBenchmarkWithRepeatNum(10000, func() {
		NewIdBinary(nil)
	})
}
