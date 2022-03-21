/*
Паттерн Строитель - это порождающий паттерн проектирования, который позволяет создавать
сложные объекты пошагово. Строитель даёт возможность использовать один и тот же код строительства
для получения разных представлений объектов.
Паттерн Строитель предлагает вынести конструирование объекта за пределы его собственного класса,
поручив это дело отдельным объектам, называемым строителями. Паттерн предлагает разбить процесс конструирования
объекта на отдельные шаги. Чтобы создать объект, нужно поочерёдно вызывать методы строителя.
Причём не нужно запускать все шаги, а только те, что нужны для производства объекта определённой конфигурации.

Плюсы:
- позволяет создавать продукты пошагово;
- позволяет использовать один и тот же код для создания различных продуктов;
- изолирует сложный код сборки продукта от его основной бизнес-логики.
Минусы:
- усложняет код программы из-за введения дополнительных классов.

В данном примере реализован строитель html документа.
*/

package main

import "fmt"

// Builder - описывает интерфейс строителя.
// В данном случае строитель содержит методы генерации различных частей html документа.
type Builder interface {
	MakeHeader(str string)
	MakeBody(str string)
	MakeFooter(str string)
}

// Director - структура, которая распоряжается строителем и отдает ему команды в нужном порядке,
// а строитель их выполняет.
type Director struct {
	builder Builder
}

// Construct - метод, генерирующий html документ.
func (d *Director) Construct() {
	d.builder.MakeHeader("Header")
	d.builder.MakeBody("Body")
	d.builder.MakeFooter("Footer")
}

// ConcreteBuilder - реализует интерфейс строителя Builder и взаимодействует со сложным объектом.
type ConcreteBuilder struct {
	product *Product
}

// MakeHeader - создает заголовок html документа.
func (b *ConcreteBuilder) MakeHeader(str string) {
	b.product.Content += "<header>" + str + "</header>"
}

// MakeBody - создает тело html документа.
func (b *ConcreteBuilder) MakeBody(str string) {
	b.product.Content += "<article>" + str + "</article>"
}

// MakeFooter - создает футер html документа.
func (b *ConcreteBuilder) MakeFooter(str string) {
	b.product.Content += "<footer>" + str + "</footer>"
}

// Product - сложный объект.
type Product struct {
	Content string
}

// Show возвращает объект.
func (p *Product) Show() string {
	return p.Content
}

func main() {
	// Экземпляр сложного объекта.
	obj := new(Product)
	// Экземпляр конкретного строителя.
	builder := ConcreteBuilder{
		product: obj,
	}
	// Экземпляр директора.
	director := Director{
		builder: &builder,
	}
	director.Construct()
	fmt.Println(obj.Show())
	// output:
	// <header>Header</header><article>Body</article><footer>Footer</footer>
}
