/*
Паттерн Цепочка обязанностей - это поведенческий паттерн проектирования,
который позволяет передавать запросы последовательно по цепочке обработчиков.
Каждый последующий обработчик решает, может ли он обработать запрос сам и стоит ли передавать запрос дальше по цепи.

Плюсы:
- уменьшает зависимость между клиентом и обработчиками.
Минусы:
- запрос может остаться никем не обработанным.
*/

package main

import "fmt"

// Handler - интерфейс, описывающий поведение обработчиков в цепочке.
type Handler interface {
	SendRequest(int) string
}

// ConcreteHandlerA - реализует конкретный обработчик A.
// Содержит ссылку на следующий обработчик.
type ConcreteHandlerA struct {
	next Handler
}

// SendRequest - обрабатывает запрос и решает, выполнять его дальше по цепочке или нет.
func (h *ConcreteHandlerA) SendRequest(message int) (result string) {
	if message == 1 {
		result = "I'm handler 1"
	} else if h.next != nil {
		result = h.next.SendRequest(message)
	}
	return
}

// ConcreteHandlerB - реализует конкретный обработчик B.
// Содержит ссылку на следующий обработчик.
type ConcreteHandlerB struct {
	next Handler
}

// SendRequest - обрабатывает запрос и решает, выполнять его дальше по цепочке или нет.
func (h *ConcreteHandlerB) SendRequest(message int) (result string) {
	if message == 2 {
		result = "I'm handler 2"
	} else if h.next != nil {
		result = h.next.SendRequest(message)
	}
	return
}

// ConcreteHandlerC - реализует конкретный обработчик C.
// Содержит ссылку на следующий обработчик.
type ConcreteHandlerC struct {
	next Handler
}

// SendRequest - обрабатывает запрос и решает, выполнять его дальше по цепочке или нет.
func (h *ConcreteHandlerC) SendRequest(message int) (result string) {
	if message == 3 {
		result = "I'm handler 3"
	} else if h.next != nil {
		result = h.next.SendRequest(message)
	}
	return
}

func main() {
	handlerA := new(ConcreteHandlerA)
	handlerB := new(ConcreteHandlerB)
	handlerC := new(ConcreteHandlerC)
	handlerA.next = handlerB
	handlerB.next = handlerC
	message := 3
	fmt.Println(handlerA.SendRequest(message))
	message = 2
	fmt.Println(handlerA.SendRequest(message))
	message = 1
	fmt.Println(handlerA.SendRequest(message))
	// output:
	// I'm handler 3
	// I'm handler 2
	// I'm handler 1
}
