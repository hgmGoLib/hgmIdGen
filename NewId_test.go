package hgmIdGen

import (
	"testing"
	"fmt"
	"time"
	"sync"
)

func TestNewId(ot *testing.T) {
	_Equal(len(NewId()),29)
	//fmt.Println(NewId())
	//fmt.Println(NewId())
	_Equal(IsValidId(NewId()),true)
	_Equal(IsValidSecureId(NewId()),false)
	t:=uint64(time.Now().UnixMilli())
	//fmt.Println(shortEncode(t,1073741824))
	last:= shortEncode(t,0)
	for i:=t+1;i<t+1000;i++{
		this:= shortEncode(i,0)
		if last>=this{
			fmt.Println(last,this)
			panic("fail")
		}
		last = this
	}
	idIncCheckFn:=func(list []uint32){
		last=""
		for _,id:=range list{
			this:= shortEncode(t,id)
			if last>=this{
				fmt.Println(last,this)
				panic("fail "+_Format(id))
			}
			last = this
		}
	}
	idIncCheckFn([]uint32{0,1,0x3F,0x3F+1,0x3FFF,0x3FFF+1,0x3FFFFF,0x3FFFFF+1,1073741822,1073741823})
	// 回环之后,不会蹦,还是能递增.
	idIncCheckFn([]uint32{1073741824,1073741825})
	last= NewId()
	for i:=0;i<1000*1000;i++{
		this:= NewId()
		isOk:=this>=last
		if isOk==false{
			fmt.Println(i,this,last)
			panic(`fail`)
		}
		last = this
		_Equal(IsValidId(this),true)
		_Equal(IsValidSecureId(this),false)
	}
	a := shortEncode(t, 0x3FFFFFFF)
	b := shortEncode(t+1, 0)
	if a >= b {
		panic("cross-time id not increasing")
	}
	{
		var wg sync.WaitGroup
		idChan := make(chan string, 10000)
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < 1000; j++ {
					idChan <- NewId()
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(idChan)
		ids := make(map[string]bool)
		for id := range idChan {
			if ids[id] {
				panic("duplicate id")
			}
			ids[id] = true
		}
	}
}

func Benchmark_NewId(ot *testing.B) {
	for i:=0;i<ot.N;i++{
		NewId()
	}
}