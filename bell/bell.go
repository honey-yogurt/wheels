package bell

import (
	"context"
	"sync/atomic"
	"time"
)

// Bell A simple timer alarm.
type Bell struct {
	alarmInterval time.Duration //报警间隔，实则是监控间隔
	alarmCount    int32         //计数器，满足业务规定的条件时候 needAlarm 一下，增加计数器
	maxAlarmCount int32         //临界值，当在一个间隔内，alarmCount >= maxAlarmCount 的时候，触发报警,业务处理报警逻辑
	alreadyAlarm  bool          //标志位，表示是否已经触发过报警
}

func NewBell(alarmInterval time.Duration, maxAlarmCount int32) *Bell {
	return &Bell{
		alarmInterval: alarmInterval,
		maxAlarmCount: maxAlarmCount,
		alarmCount:    0,
		alreadyAlarm:  false,
	}
}

func (b *Bell) start(ctx context.Context) {
	ticker := time.NewTicker(b.alarmInterval)
	for {
		select {
		case <-ticker.C:
			//当定时器触发时，将报警次数重置为零，并将已触发报警的标志位设为 false
			atomic.StoreInt32(&b.alarmCount, 0)
			b.alreadyAlarm = false
		case <-ctx.Done():
			return
		}
	}
}

// 检查是否需要触发报警
func (b *Bell) needAlarm() bool {
	n := atomic.AddInt32(&b.alarmCount, 1)
	// beyond maxAlarmCount and did not alarm in time space
	if n > b.maxAlarmCount && !b.alreadyAlarm {
		b.alreadyAlarm = true
		return true
	}
	return false
}
