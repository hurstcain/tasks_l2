/*
Паттерн Стратегия используется, когда есть семейство некоторых схожих алгоритмов, которые часто
изменяются или расширяются. Тогда, согласно паттерну Стратегия, каждый из этих алгоритмов помещается в свой
собственный класс (их можно назвать конкретными стратегиями), эти классы реализуют один и тот же интерфейс.
Затем некоторый основной класс, вместо того, чтобы самому реализовывать алгоритм, будет ссылаться на один
из классов-стратегий и делегировать реализацию алгоритма этому классу.

Плюсы:
- изолирует код алгоритмов от основного класса, что позволяет не трогать код основного класса при изменении
или добавлении новых алгоритмов;
- так как классы-стратегии реализуют один интерфейс, с помощью сеттера можно изменять
используемый в данный момент алгоритм;
- вместо наследования используется композиция.
Минусы:
- из-за дополнительных классов усложняется код программы;
- клиент должен знать различия между стратегиями, чтобы использовать их в коде.
*/

package pattern

import "fmt"

// Operator - общий интерфейс для всех стратегий.
// Данный интерфейс описывает структуру с методом, который принимает в качестве аргумента
// два целых числа, выполняет некоторые операции с этими числами и возвращает результат этих операций.
type Operator interface {
	Execute(int, int) int
}

// Addition - описывает алгоритм суммирования двух чисел.
type Addition struct{}

// Execute - сложение двух чисел a и b.
func (add Addition) Execute(a, b int) int {
	return a + b
}

// Multiplication - описывает алгоритм умножения двух чисел.
type Multiplication struct{}

// Execute - умножение двух чисел a и b.
func (n Multiplication) Execute(a, b int) int {
	return a * b
}

// Subtraction - описывает алгоритм вычитания двух чисел.
type Subtraction struct{}

// Execute - вычитание двух чисел a и b.
func (s Subtraction) Execute(a, b int) int {
	return a - b
}

// Division - описывает алгоритм деления двух чисел.
type Division struct{}

// Execute - деление двух чисел a и b.
func (d Division) Execute(a, b int) int {
	return a / b
}

// Operation - структура, которая делегирует реализацию алгоритма некоторой
// стратегии op (т.е. одной из структур, реализующих интерфейс Operator).
type Operation struct {
	op Operator
}

// NewOperation - конструктор структуры Operation.
func NewOperation(op Operator) Operation {
	return Operation{
		op: op,
	}
}

// SetOperation - устанавливает новый алгоритм.
func (o *Operation) SetOperation(op Operator) {
	o.op = op
}

// DoOperation - выполняет некоторую операцию, алгоритм которой определен в Operation.op.
func (o Operation) DoOperation(a, b int) int {
	return o.op.Execute(a, b)
}

// Client - код клиента.
func Client() {
	a, b := 20, 10

	o := NewOperation(Addition{})
	fmt.Printf("%d + %d = %d\n", a, b, o.DoOperation(a, b))
	o.SetOperation(Multiplication{})
	fmt.Printf("%d * %d = %d\n", a, b, o.DoOperation(a, b))
	o.SetOperation(Subtraction{})
	fmt.Printf("%d - %d = %d\n", a, b, o.DoOperation(a, b))
	o.SetOperation(Division{})
	fmt.Printf("%d / %d = %d\n", a, b, o.DoOperation(a, b))
	// output:
	// 20 + 10 = 30
	// 20 * 10 = 200
	// 20 - 10 = 10
	// 20 / 10 = 2
}
