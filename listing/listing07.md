Что выведет программа? Объяснить вывод программы.

```go
package main
import (
	"fmt"
	"math/rand"
	"time"
)
func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}
func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}
func main() {
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
Вывод:
Сначала в случайном порядке выводятся числа 1 3 5 7 2 4 6 8, а потом бесконечно выводятся нули.

Это происходит, так как после закрытия двух каналов в функции merge происходит бесконечное чтение
из закрытых каналов значений по умолчанию. Чтобы этого не происходило, нужно добавить проверку
на то, закрыт канал или нет. 
```