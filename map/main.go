package main

import (
	"log"

	"sync"
	"time"
)

var (
	max    = 5000
	conMap = NewRWMap(make(map[int]*strWithLock, max))
)

type strWithLock struct {
	content int
	mutex   sync.Mutex
}

func (s *strWithLock) Update() {
	s.mutex.Lock()
	if s.content%2 == 0 {
		conMap.Delete(s.content)
	}
	s.mutex.Unlock()
}

func main() {
	log.Printf("start")
	start := time.Now()
	defer func() {
		log.Printf("cost: %v", time.Since(start))
	}()

	// 写map
	for i := 0; i < max; i++ {
		go func(index int) {
			newStrLock := &strWithLock{content: index}
			conMap.Set(index, newStrLock)
		}(i)
	}
	time.Sleep(time.Second * 1)

	// 修改map中的元素
	for i := 0; i < 20; i++ {
		go func() {
			for i := 0; i < max; i++ {
				go func(index int) {
					if v, ok := conMap.Get(index); ok {
						v.Update()
					}
				}(i)
			}
		}()
	}

	// 读map
	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < max; i++ {
				go func(index int) {
					_, _ = conMap.Get(index)
				}(i)
			}
		}()
	}
	// 写map
	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < max; i++ {
				go func(index int) {
					if _, ok := conMap.Get(index); !ok {
						newStrLock := &strWithLock{content: index}
						conMap.Set(index, newStrLock)
					}
				}(i)
			}
		}()
	}
	time.Sleep(time.Second * 7)
	log.Printf("len of map: %v", conMap.Len())
}
