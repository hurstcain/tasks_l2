/*
Паттерн Команда - поведенческий паттерн проектирования, который превращает запросы в объекты,
позволяя передавать их как аргументы при вызове методов, ставить в очередь,
логировать их, а также поддерживать отмену операций.
Данный паттерн отделяет объект, инициирующий операцию, от объекта,
который знает, как ее выполнить. Единственное, что должен знать инициатор, это как отправить команду.
Команда - описывает общий интерфейс для всех конкретных команд.
Получатель - содержит бизнес-логику программы, умеет выполнять все виды операций, связанных с выполнением запроса.
Инициатор - выполняет команды, поддерживает добавление, а также удаление команд на выполнение.
Связан с одной или несколькими командами.

Плюсы:
- объекты, вызывающие операции, и объекты, которые их выполняют, изолированы друг от друга;
- реализует повтор, отмену и отложенный запуск операций;
- позволяет собирать сложные команды из простых.
Минусы:
- усложняет код программы из-за добавления дополнительных классов.

Примеры использования:
Используется, когда нужно откладывать выполнение команд, выстраивать их в очереди,
а также хранить историю и делать отменять выполнение команд.
*/

package pattern

import "fmt"

// Получатель.

// Receiver - получатель запроса, обрабатывающий команды.
type Receiver struct{}

// On - реализует кнопку включения.
func (r *Receiver) On() string {
	return "Turn On"
}

// Off - реализует кнопку выключения.
func (r *Receiver) Off() string {
	return "Turn Off"
}

// Интерфейс команды.

// Command - интерфейс команды, объявляет метод Execute для выполнения команд.
type Command interface {
	Execute() string
}

// Конкретные команды.

// OnCommand - команда включения, реализует интерфейс Command.
type OnCommand struct {
	receiver *Receiver
}

// Execute - выполняет включение.
func (c *OnCommand) Execute() string {
	return c.receiver.On()
}

// OffCommand команда выключения, реализует интерфейс Command.
type OffCommand struct {
	receiver *Receiver
}

// Execute - выполняет выключение.
func (c *OffCommand) Execute() string {
	return c.receiver.Off()
}

// Инициатор.

// Invoker - структура инициатор запроса.
type Invoker struct {
	commands []Command
}

// AddCommand - добавляет команду в список команд на исполнение.
func (i *Invoker) AddCommand(command Command) {
	i.commands = append(i.commands, command)
}

// RemoveCommand - удаляет последнюю команду из списка на исполнение.
func (i *Invoker) RemoveCommand() {
	if len(i.commands) != 0 {
		i.commands = i.commands[:len(i.commands)-1]
	}
}

// Execute - выполняет команды.
func (i *Invoker) Execute() string {
	var result string

	for _, command := range i.commands {
		result += command.Execute() + "\n"
	}

	return result
}

// CommandClient - код клиента.
func CommandClient() {
	// Экземпляр получателя.
	receiver := new(Receiver)
	// Экземпляр исполнителя команды включения.
	on := OnCommand{receiver: receiver}
	// Экземпляр исполнителя команды выключения.
	off := OffCommand{receiver: receiver}
	// Экземпляр инициатора команд.
	invoker := new(Invoker)
	// Добавляем команды.
	invoker.AddCommand(&on)
	invoker.AddCommand(&off)
	invoker.AddCommand(&off)
	invoker.AddCommand(&off)
	// Удаляем последнюю команду.
	invoker.RemoveCommand()
	// Вывод результата исполнения команд на экран.
	fmt.Print(invoker.Execute())
	// Output:
	// Turn On
	// Turn Off
	// Turn Off
}
