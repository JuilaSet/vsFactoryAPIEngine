package util

type Signal chan interface{}

// 信号量的初始值
func NewSignal(i int) Signal {
	return make(Signal, i)
}

// 尝试占取资源
func P(s Signal) {
	s <- 0
}

// 释放资源
func V(s Signal) {
	<-s
}
