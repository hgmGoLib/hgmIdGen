package hgmIdGen

import (
	"testing"
	"time"
)

func TestParseTimeFromId(ot *testing.T) {
	expect := time.UnixMilli(1741257015123)
	id := shortEncode(uint64(expect.UnixMilli()), 123)
	got := ParseTimeFromId(id)
	_Equal(got, expect)
}

func TestParseTimeFromSecureId(ot *testing.T) {
	expect := time.UnixMilli(1741257015123)
	id := secureEncode(uint64(expect.UnixMilli()), 456)
	got := ParseTimeFromSecureId(id)
	_Equal(got, expect)
}

func TestParseTimeFromIdFail(ot *testing.T) {
	_Equal(ParseTimeFromId(""), time.Time{})
	_Equal(ParseTimeFromId("17fatw9hnw1y8ge6emw3et3fmtu9!"), time.Time{})
}

func TestParseTimeFromSecureIdFail(ot *testing.T) {
	_Equal(ParseTimeFromSecureId(""), time.Time{})
	_Equal(ParseTimeFromSecureId("17fatweptd1swfgg11ekgtx5j459wnf112v7jp81w!"), time.Time{})
}

func TestMinIdAtTime(ot *testing.T) {
	expect := time.UnixMilli(1741257015123)
	minId := MinIdAtTime(expect)
	_Equal(len(minId), 29)
	_Equal(IsValidId(minId), true)
	_Equal(ParseTimeFromId(minId), expect)
	_Equal(minId < shortEncode(uint64(expect.UnixMilli()), 1), true)
	_Equal(minId < MinIdAtTime(expect.Add(time.Millisecond)), true)
}

func TestMinSecureIdAtTime(ot *testing.T) {
	expect := time.UnixMilli(1741257015123)
	minId := MinSecureIdAtTime(expect)
	_Equal(len(minId), 42)
	_Equal(IsValidSecureId(minId), true)
	_Equal(ParseTimeFromSecureId(minId), expect)
	_Equal(minId < secureEncode(uint64(expect.UnixMilli()), 1), true)
	_Equal(minId < MinSecureIdAtTime(expect.Add(time.Millisecond)), true)
}
