/*
Паттерн Состояние - это поведенческий паттерн проектирования, который позволяет объектам менять поведение
в зависимости от своего состояния. Извне создаётся впечатление, что изменился класс объекта.
Данный паттерн предлагает создавать отдельные классы для каждого состояния, в котором может пребывать объект,
а затем вынести туда поведения, соответствующие этим состояниям.
Вместо того чтобы хранить код всех состояний, первоначальный объект, называемый контекстом,
будет содержать ссылку на один из объектов-состояний и делегировать ему работу, зависящую от состояния.
Благодаря тому, что объекты состояний будут иметь общий интерфейс, контекст сможет делегировать работу состоянию,
не привязываясь к его классу. Поведение контекста можно будет изменить в любой момент,
подключив к нему другой объект-состояние.

Плюсы:
- избавляет от множества больших условных операторов машины состояний;
- концентрирует в одном месте код, связанный с определённым состоянием;
- упрощает код контекста.
Минусы:
- может неоправданно усложнить код, если состояний мало и они редко меняются.
*/

package pattern

import "fmt"

// State - описывает общий интерфейс для всех конкретных состояний.
// В данном случае все состояния реализуют методы FirstHandler и SecondHandler.
type State interface {
	FirstHandler()
	SecondHandler()
}

// Context - структура, в которой хранится ссылка на текущее состояние объекта (currentState).
// В зависимости от текущего состояния, поведение объекта меняется.
type Context struct {
	currentState State
}

// SetState - Устанавливает состояние объекта.
func (c *Context) SetState(state State) {
	c.currentState = state
}

// FirstHandler - первый обработчик, реализация зависит от текущего состояния.
func (c *Context) FirstHandler() {
	c.currentState.FirstHandler()
}

// SecondHandler - второй обработчик, реализация зависит от текущего состояния.
func (c *Context) SecondHandler() {
	c.currentState.SecondHandler()
}

// ConcreteStateA - конкретное состояние A.
type ConcreteStateA struct {
	ctx *Context
}

// FirstHandler - первый обработчик, выводит некоторую информацию на экран.
func (c *ConcreteStateA) FirstHandler() {
	fmt.Println("First handler of State A")
}

// SecondHandler - второй обработчик, выводит некоторую информацию на экран.
// Затем меняет состояние объекта.
func (c *ConcreteStateA) SecondHandler() {
	fmt.Println("Second handler of State A")
	fmt.Println("Change object state")
	c.ctx.SetState(&ConcreteStateB{ctx: c.ctx})
}

// ConcreteStateB - - конкретное состояние B.
type ConcreteStateB struct {
	ctx *Context
}

// FirstHandler - первый обработчик, выводит некоторую информацию на экран.
func (c *ConcreteStateB) FirstHandler() {
	fmt.Println("First handler of State B")
}

// SecondHandler - второй обработчик, выводит некоторую информацию на экран.
// Затем меняет состояние объекта.
func (c *ConcreteStateB) SecondHandler() {
	fmt.Println("Second handler of State B")
	fmt.Println("Change object state")
	c.ctx.SetState(&ConcreteStateA{ctx: c.ctx})
}

// StateClient - код клиента.
func StateClient() {
	ctx := new(Context)
	ctx.SetState(&ConcreteStateA{ctx: ctx})
	ctx.FirstHandler()
	ctx.SecondHandler()
	ctx.FirstHandler()
	ctx.SecondHandler()
	ctx.FirstHandler()
	// output:
	// First handler of State A
	// Second handler of State A
	// Change object state
	// First handler of State B
	// Second handler of State B
	// Change object state
	// First handler of State A
}
