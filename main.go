package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: pomodoro N")
		os.Exit(1)
	}
	var (
		n   int
		err error
	)
	if n, err = strconv.Atoi(os.Args[1]); err != nil {
		fmt.Println("Usage: pomodoro N, N must be a number")
		os.Exit(1)
	}
	p := NewPomodoro(n)
	p.Start()
}

func countdown(d time.Duration) <-chan int {
	tick := make(chan int, 1)
	go func() {
		c := int(d / time.Second)
		for i := c; i >= 0; i-- {
			tick <- i
			<-time.After(time.Second)
		}
		close(tick)
	}()
	return tick
}

const (
	S_WAITING_START = 0
	S_WORK          = 1
	S_BREAK         = 2
	S_END           = 3
	T_WK            = 25 * time.Minute
	T_BK            = 5 * time.Minute
)

type pomodoro struct {
	status int
	n      int
}

func NewPomodoro(n int) pomodoro {
	return pomodoro{S_WAITING_START, n}
}

func (p *pomodoro) Start() {
	for i := 0; i < p.n*2-1; i++ {
		switch p.status {
		case S_WAITING_START, S_BREAK:
			p.wk()
			p.status = S_WORK
		case S_WORK:
			p.bk()
			p.status = S_BREAK
		}
	}
	p.st()
}
func (p pomodoro) wk() {
	ct := countdown(T_WK)
	for r := range ct {
		displayCountdown(r)
	}
}
func (p pomodoro) bk() {
	ct := countdown(T_BK)
	for rs := range ct {
		displayCountdown(rs)
	}
}

func (p pomodoro) st() {
	displayEnd()
}

func displayCountdown(seconds int) {
	h := seconds / 3600
	m := seconds % 3600 / 60
	s := seconds % 3600 % 60
	fmt.Printf("\r%.2d:%.2d:%.2d", h, m, s)
}

func displayEnd() {
	fmt.Println()
}
