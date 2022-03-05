package main

import (
	"fmt"
	"sync"
	"time"
)

// Функция, которая объединяет один или более done-каналов в single-канал,
// если хотя бы один из этих done-каналов окажется закрыт.
func or(channels ...<-chan interface{}) <-chan interface{} {
	// Структура single-канала.
	// Состоит из самого канала ch.
	// Флага closed, который определяет закрыт канал или нет.
	// И встроенных методов структуры sync.Mutex.
	singleCh := struct {
		ch     chan interface{}
		closed bool
		sync.Mutex
	}{
		ch:     make(chan interface{}),
		closed: false,
	}

	// Группа горутин, завершения которых будет ожидать функция.
	wg := sync.WaitGroup{}

	// Запуск горутин, каждая из которых будет ожидать закрытия одного из каналов channels и канала singleCh.
	for _, val := range channels {
		wg.Add(1)
		go func(done <-chan interface{}) {
			defer wg.Done()
			for {
				select {
				case _, ok := <-done:
					// Если канал не закрыт, переходим на следующую итерацию цикла.
					if ok {
						continue
					}
					// Если канал done (один из каналов channels) оказался закрыт, то нужно закрыть канал singleCh.
					// Сначала блокируем доступ к каналу.
					singleCh.Lock()
					// Если канал уже закрыт, то горутина завершает работу.
					if singleCh.closed {
						singleCh.Unlock()
						return
					}
					// Если канал не закрыт, то закрываем канал и присваиваем флагу закрытия канала true.
					close(singleCh.ch)
					singleCh.closed = true
					singleCh.Unlock()
				case <-singleCh.ch:
					// Если канал singleCh оказался закрыт раньше, чем закрылся канал done,
					// то горутина завершается.
					return
				}
			}
		}(val)
	}

	// Ожидание завершения горутин.
	wg.Wait()

	return singleCh.ch
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-or(
		sig(10*time.Second),
		sig(6*time.Second),
		sig(11*time.Second),
		sig(6*time.Second),
		sig(12*time.Second),
	)

	fmt.Printf("done after %v", time.Since(start))
	// output:
	// done after 6.0008733s
}
