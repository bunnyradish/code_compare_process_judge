package controlroutine

import "judge/zapconf"

/*
控制协程数量
*/

type ChanRoutine struct {
	ch chan struct{}
}

func NewChanRoutine(n int) *ChanRoutine {//相当于channel就是一个控制goroutine数量的队列，把channel的缓存量开为8，然后后面的add和del就是生产者消费者模式
	cr := &ChanRoutine{make(chan struct{}, n)}
	return cr
}

func (cr *ChanRoutine) AddGoRoutine() {
	cr.ch <- struct{}{}
}

func (cr *ChanRoutine) DelGoRoutine() {
	select {
	case <- cr.ch:
	default:
		zapconf.GetInfoLog().Info("del chan err")
		panic("del chan err")
	}
}

