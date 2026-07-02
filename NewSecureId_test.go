package hgmIdGen

import (
	"testing"
	"fmt"
	"time"
)

func TestNewSecureId(ot *testing.T) {
	_Equal(len(NewSecureId()),42)
	fmt.Println(NewSecureId())
	fmt.Println(NewSecureId())
	_Equal(IsValidId(NewSecureId()),false)
	_Equal(IsValidSecureId(NewSecureId()),true)
	t:=uint64(time.Now().UnixMilli())
	last:= secureEncode(t,0)
	for i:=t+1;i<t+1000;i++{
		this:= secureEncode(i,0)
		if last>=this{
			fmt.Println(last,this)
			panic("fail")
		}
		last = this
	}
	last= NewSecureId()
	for i:=0;i<1000*1000;i++{
		this:= NewSecureId()
		isOk:=this>=last
		if isOk==false{
			fmt.Println(i,this,last)
			panic(`fail`)
		}
		last = this
		_Equal(IsValidId(this),false)
		_Equal(IsValidSecureId(this),true)
	}
	// Benchmark [175.5ns/op] [5.435MBop/s] duration:[175.5us] allocNum:[1.0000B/op] allocSize:[32.000B/op] 
	//BenchmarkWithRepeatNum(1000, func() {
	//	NewSecureId()
	//})
}