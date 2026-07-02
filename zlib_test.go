package hgmIdGen

// 本文件从 hgmLib/hgmTest 与 hgmLib/hgmStrconv 复制而来, 仅供本库自动测试使用.
// 使本库独立, 不依赖 hgmLib. 相等判断语义与 hgmTest.Equal 一致; 仅失败时的 diff 展示改用 fmt 简化.

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 不相等 就打出 messageWhenFail的消息然后 panic (停止测试)
func _Equal[T any](a T, b T, messageWhenFail ...any) {
	if reflect.DeepEqual(a, b) {
		return
	}
	msg := fmt.Sprintf("hgmTest.Equal fail type:%s\na: %#v\nb: %#v\n", _TypeDebugAny(a), a, b)
	if len(messageWhenFail) > 0 {
		msg += "messageWhenFail: " + fmt.Sprintln(messageWhenFail...)
	}
	os.Stderr.Write([]byte(msg))
	panic("hgmTest fail")
}

func _TypeDebugAny(obj any) string {
	if obj == nil {
		return "nil"
	}
	return reflect.TypeOf(obj).String()
}

// 进行一次benchmark,同一个进程只允许同时运行一个.
func _hgmBenchmark(fn func()) {
	isOk := gBenchmarkCtx.isCalling.CompareAndSwap(false, true)
	if isOk == false {
		panic(`[hgmTest.Benchmark] reentry, please do not call Benchmark in Benchmark`)
	}
	defer gBenchmarkCtx.isCalling.Store(false)
	gBenchmarkCtx.runLocker.Lock()
	defer gBenchmarkCtx.runLocker.Unlock()
	gBenchmarkCtx.dataLocker.Lock()
	if gBenchmarkCtx.num == 0 {
		gBenchmarkCtx.num = 1
	}
	gBenchmarkCtx.startCgoCall = runtime.NumCgoCall()
	gBenchmarkCtx.dataLocker.Unlock()
	var memstats1 runtime.MemStats
	var memstats2 runtime.MemStats
	runtime.ReadMemStats(&memstats1)
	startTime := time.Now()
	fn()
	dur := time.Since(startTime)
	runtime.ReadMemStats(&memstats2)
	gBenchmarkCtx.dataLocker.Lock()
	gBenchmarkCtx.endCgoCall = runtime.NumCgoCall()
	gBenchmarkCtx.dur = dur
	gBenchmarkCtx.allocNum = memstats2.Mallocs - memstats1.Mallocs
	gBenchmarkCtx.allocSize = memstats2.TotalAlloc - memstats1.TotalAlloc
	resultS := _benchmarkResultString__NOLOCK()
	gBenchmarkCtx.name = ""
	gBenchmarkCtx.sizePerRun = 0
	gBenchmarkCtx.num = 0
	gBenchmarkCtx.dataLocker.Unlock()
	_, _ = os.Stdout.WriteString(resultS)
	_, _ = os.Stdout.WriteString("\n")
}

func _hgmBenchmarkWithRepeatNum(num int, fn func()) {
	_hgmBenchmark(func() {
		_hgmBenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			fn()
		}
	})
}

type benchmarkCtx struct {
	num             int
	sizePerRun      int
	dur             time.Duration
	name            string
	namePadding     int
	allocNum        uint64
	allocSize       uint64
	runLocker       sync.Mutex
	dataLocker      sync.Mutex
	startCgoCall    int64
	endCgoCall      int64
	namePaddingSize int
	isCalling       atomic.Bool
}

var gBenchmarkCtx benchmarkCtx

func _hgmBenchmarkSetNum(num int) {
	gBenchmarkCtx.dataLocker.Lock()
	gBenchmarkCtx.num = num
	gBenchmarkCtx.dataLocker.Unlock()
}

func _hgmBenchmarkSetName(name string) {
	gBenchmarkCtx.dataLocker.Lock()
	gBenchmarkCtx.name = name
	gBenchmarkCtx.dataLocker.Unlock()
}

func _benchmarkResultString__NOLOCK() string {
	buf := bytes.Buffer{}
	buf.WriteString("Benchmark ")
	if gBenchmarkCtx.name != "" {
		buf.WriteString(gBenchmarkCtx.name)
		if gBenchmarkCtx.namePadding <= 1 {
			gBenchmarkCtx.namePadding = 8
		}
		if len(gBenchmarkCtx.name) <= gBenchmarkCtx.namePadding-1 {
			buf.WriteString(strings.Repeat(" ", gBenchmarkCtx.namePadding-1-len(gBenchmarkCtx.name)))
		}
		buf.WriteString(" ")
	}
	buf.WriteString("[")
	buf.WriteString(_durationFormatFloat64Ns(float64(gBenchmarkCtx.dur) / float64(gBenchmarkCtx.num)))
	buf.WriteString("/op] ")

	buf.WriteString("[")
	buf.WriteString(_gbFromFloat64(float64(gBenchmarkCtx.num) / float64(gBenchmarkCtx.dur) * 1e9))
	buf.WriteString("op/s] ")

	buf.WriteString("duration:[")
	buf.WriteString(_durationFormatFloat64Ns(float64(gBenchmarkCtx.dur)))
	buf.WriteString("] ")

	buf.WriteString("allocNum:[")
	buf.WriteString(_gbFromFloat64(float64(gBenchmarkCtx.allocNum) / float64(gBenchmarkCtx.num)))
	buf.WriteString("/op] ")

	buf.WriteString("allocSize:[")
	buf.WriteString(_gbFromFloat64(float64(gBenchmarkCtx.allocSize) / float64(gBenchmarkCtx.num)))
	buf.WriteString("/op] ")

	if gBenchmarkCtx.sizePerRun > 0 {
		buf.WriteString("bandwidth:[")
		buf.WriteString(_gbFromFloat64(float64(gBenchmarkCtx.num*gBenchmarkCtx.sizePerRun) / float64(gBenchmarkCtx.dur) * 1e9))
		buf.WriteString("/s] ")
	}
	if gBenchmarkCtx.endCgoCall-gBenchmarkCtx.startCgoCall > 0 {
		buf.WriteString("cgoNum:[")
		buf.WriteString(_gbFromFloat64(float64(gBenchmarkCtx.endCgoCall-gBenchmarkCtx.startCgoCall) / float64(gBenchmarkCtx.num)))
		buf.WriteString("/op] ")
	}
	return buf.String()
}

// 使用float64的ns作为时间单位(一般是做了很多事情之后的平均值)
func _durationFormatFloat64Ns(dur float64) string {
	const day = 24 * time.Hour
	const year = 365 * day
	if (dur >= float64(year)) || (dur <= float64(-year)) {
		return _formatFloat64ToFInLen(float64(dur)/float64(year), 6) + "y"
	} else if dur >= float64(day) || dur <= float64(-day) {
		return _formatFloat64ToFInLen(float64(dur)/float64(day), 6) + "d"
	} else if dur >= float64(time.Hour) || dur < float64(-time.Hour) {
		return _formatFloat64ToFInLen(float64(dur)/float64(time.Hour), 6) + "h"
	} else if dur >= float64(time.Minute) || dur <= float64(-time.Minute) {
		return _formatFloat64ToFInLen(float64(dur)/float64(time.Minute), 6) + "m"
	} else if dur >= float64(time.Second) || dur <= float64(-time.Second) {
		return _formatFloat64ToFInLen(float64(dur)/float64(time.Second), 6) + "s"
	} else if dur >= float64(time.Millisecond) || dur <= float64(-time.Millisecond) {
		return _formatFloat64ToFInLen(float64(dur)/float64(time.Millisecond), 5) + "ms"
	} else if dur >= float64(time.Microsecond) || dur <= float64(-time.Microsecond) {
		return _formatFloat64ToFInLen(float64(dur)/float64(time.Microsecond), 5) + "us"
	} else {
		return _formatFloat64ToFInLen(dur, 5) + "ns"
	}
}

func _gbFromFloat64(byteNum float64) string {
	if byteNum >= 1e15 || byteNum <= -1e15 {
		return _formatFloat64ToFInLen(byteNum/(1024*1024*1024*1024*1024), 5) + "PB"
	}
	if byteNum >= 1e12 || byteNum <= -1e12 {
		return _formatFloat64ToFInLen(byteNum/(1024*1024*1024*1024), 5) + "TB"
	}
	if byteNum >= 1e9 || byteNum <= -1e9 {
		return _formatFloat64ToFInLen(byteNum/(1024*1024*1024), 5) + "GB"
	}
	if byteNum >= 1e6 || byteNum <= -1e6 {
		return _formatFloat64ToFInLen(byteNum/(1024*1024), 5) + "MB"
	}
	if byteNum >= 1e3 || byteNum <= -1e3 {
		return _formatFloat64ToFInLen(byteNum/(1024), 5) + "KB"
	}
	return _formatFloat64ToFInLen(byteNum, 6) + "B"
}

func _formatFloat64ToFInLen(f float64, showLen int) string {
	s1 := strconv.FormatFloat(f, 'f', 0, 64)
	if len(s1)+1 >= showLen {
		if len(s1) == showLen {
			return s1
		} else {
			return "0" + s1
		}
	}
	return strconv.FormatFloat(f, 'f', showLen-len(s1)-1, 64)
}

// 一种简单数据的"常识性" format (来自 hgmStrconv.Format)
func _Format[T string | bool | uint64 | uint32 | uint16 | uint8 | uint | int64 | int32 | int16 | int8 | int | float64 | float32](obj T) string {
	switch objI := any(obj).(type) {
	case uint64:
		return strconv.FormatUint(objI, 10)
	case uint32:
		return strconv.FormatUint(uint64(objI), 10)
	case uint16:
		return strconv.FormatUint(uint64(objI), 10)
	case uint8:
		return strconv.FormatUint(uint64(objI), 10)
	case uint:
		return strconv.FormatUint(uint64(objI), 10)
	case int64:
		return strconv.FormatInt(objI, 10)
	case int32:
		return strconv.FormatInt(int64(objI), 10)
	case int16:
		return strconv.FormatInt(int64(objI), 10)
	case int8:
		return strconv.FormatInt(int64(objI), 10)
	case int:
		return strconv.FormatInt(int64(objI), 10)
	case string:
		return objI
	case float64:
		return strconv.FormatFloat(objI, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(objI), 'f', -1, 32)
	case bool:
		return strconv.FormatBool(objI)
	default:
		return "[Format] not support type " + reflect.TypeOf(objI).String()
	}
}
