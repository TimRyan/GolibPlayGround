package common

import (
	"sync"
	"time"
)

//Golib时间类
type SysTime struct {
	time.Time
	now time.Time //非自然时间
}

//实例化时间
func newSysTime(t time.Time) *SysTime {

	sysTime := new(SysTime)
	sysTime.now = t
	return sysTime
}

//单例模式
var sysTimeInstance *SysTime
var sysTimeOnce sync.Once

func GetSysTime() *SysTime {

	sysTimeOnce.Do(func() {
		sysTimeInstance = newSysTime(time.Now())
	})
	return sysTimeInstance
}

//返回当前时间
func (p *SysTime) Now(isNatual bool) time.Time {
	if isNatual {
		return time.Now()
	} else {
		return p.now
	}
}

//时间同步
func (p *SysTime) sync(now time.Time) {
	p.now = now
}
