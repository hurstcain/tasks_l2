package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Коды цветов текста, выводимого в консоли.
const (
	redColor     = "\033[31m"
	greenColor   = "\033[32m"
	regularColor = "\033[0m"
)

// Структура, хранящая информацию о строках, записанных в stdin до совпадения.
// index - индекс строки.
// s - строка.
type beforeData struct {
	index int
	s     string
}

// Конструктор структуры beforeData.
func newBeforeData(index int, s string) beforeData {
	return beforeData{
		index: index,
		s:     s,
	}
}

// Структура, описывающая канал, в который записываются строки, переданные в stdin до совпадения.
// ch - канал, куда записываются данные.
// chCapacity - емкость канала.
type beforeDataChannel struct {
	ch         chan beforeData
	chCapacity int
}

// Конструктор структуры beforeDataChannel.
func newBeforeDataChannel(chanCapacity int) beforeDataChannel {
	return beforeDataChannel{
		ch:         make(chan beforeData, chanCapacity),
		chCapacity: chanCapacity,
	}
}

// WriteBeforeData - записывает данные в канал.
func (c *beforeDataChannel) WriteBeforeData(index int, s string) {
	// Если длина канал сравнялась с емкостью, то освобождаем канал от одной записи,
	// чтобы добавить актуальные данные.
	if len(c.ch) == c.chCapacity {
		<-c.ch
	}
	// Записываем в канал новые данные.
	c.ch <- newBeforeData(index, s)
}

// PrintBeforeData - выводит на экран данные, которые были записаны в stdin до совпадения.
// lineNum - флаг, который определяет, нужно ли выводить индекс строки или нет.
func (c *beforeDataChannel) PrintBeforeData(lineNum bool) {
	// Длина канала.
	chLen := len(c.ch)

	// Читаем все данные, которые были записаны в канал, и выводим их построчно на экран.
	for i := 0; i < chLen; i++ {
		data := <-c.ch
		if lineNum {
			fmt.Printf("%s%d: %s%s\n", greenColor, data.index, regularColor, data.s)
		} else {
			fmt.Print(lineNum, regularColor, data.s, "\n")
		}
	}
}

// Структура, в которой хранятся ключи, переданные при запуске программы.
type grepFlags struct {
	// Паттерн для поиска совпадений.
	Pattern string
	// Название файла.
	FileName string
	// Количество строк после совпадения.
	After int
	// Количество строк до совпадения.
	Before int
	// Вывод количества строк, которые соответствуют паттерну.
	Count bool
	// Флаг игнорирования регистра.
	IgnoreCase bool
	// Флаг, обозначающий исключение вместо совпадения.
	Invert bool
	// Флаг, определяющий поиск точного совпадения с введенной строкой.
	Fixed bool
	// Флаг, определяющий печатать номер строки или нет.
	LineNum bool
}

// Парсит флаги и инициализирует структуру grepFlags.
func newGrepFlags() grepFlags {
	// Количество строк после совпадения.
	var after int
	// Количество строк до совпадения.
	var before int
	flag.Func("A", "Печатать n строк после совпадения", func(aValue string) error {
		aInt, err := strconv.Atoi(aValue)
		if err != nil {
			return err
		}
		if aInt > 0 {
			after = aInt
		}
		return nil
	})
	flag.Func("B", "Печатать n строк до совпадения", func(bValue string) error {
		bInt, err := strconv.Atoi(bValue)
		if err != nil {
			return err
		}
		if bInt > 0 {
			before = bInt
		}
		return nil
	})
	flag.Func("C", "Печатать n строк вокруг совпадения", func(cValue string) error {
		cInt, err := strconv.Atoi(cValue)
		if err != nil {
			return err
		}
		if cInt > 0 {
			after = cInt
			before = cInt
		}
		return nil
	})
	count := flag.Bool("c", false, "Вывод количества строк")
	ignoreCase := flag.Bool("i", false, "Игнорировать регистр")
	invert := flag.Bool("v", false, "Вместо совпадения исключать")
	fixed := flag.Bool("F", false, "Точное совпадение со строкой, не паттерн")
	lineNum := flag.Bool("n", false, "Напечатать номер строки")
	flag.Parse()

	// Строка-паттерн для поиска совпадений. Паттерном для поиска считается первый аргумент не флаг.
	pattern := flag.Arg(0)
	if pattern == "" {
		log.Fatalln("Не был введен паттерн для поиска совпадений.")
	}
	// Название файла. Поиск совпадений в файле опционален. Названием файла считается второй аргумент.
	fileName := flag.Arg(1)

	return grepFlags{
		Pattern:    pattern,
		FileName:   fileName,
		After:      after,
		Before:     before,
		Count:      *count,
		IgnoreCase: *ignoreCase,
		Invert:     *invert,
		Fixed:      *fixed,
		LineNum:    *lineNum,
	}
}

// Grep - структура, реализующая функционал линуксовой утилиты grep.
type Grep struct {
	// Слайс, содержащий строки исходного файла. Используется, когда поиск совпадений осуществляется в файле.
	fileContent []string
	// Количество совпадений.
	matchesCount int
	flags        grepFlags
}

// NewGrep - конструктор структуры Grep.
// Инициализирует переданные при запуске программы ключи и считывает данные из исходного файла,
// если он указан.
func NewGrep() Grep {
	flags := newGrepFlags()

	// Если название файла - пустая строка, то поиск будет происходить во время ввода данных в stdin.
	if flags.FileName == "" {
		return Grep{
			matchesCount: 0,
			flags:        flags,
		}
	}

	// Открываем файл с именем fileName.
	f, err := os.Open(flags.FileName)
	if err != nil {
		log.Fatalf("Не удалось открыть исходный файл. Ошибка: %v\n", err)
	}
	// После завершения работы функции файл закрывается.
	defer f.Close()

	// Слайс, в который будут записаны строки из файла.
	fileContent := make([]string, 0)

	// Построчно читаем данные из файла.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	return Grep{
		fileContent:  fileContent,
		matchesCount: 0,
		flags:        flags,
	}
}

// Search - определяет, где будет осуществляться поиск совпадений (в файле или в stdin),
// и вызывает соответствующий метод.
func (g Grep) Search() {
	// Если название файла - пустая строка, то поиск происходит в stdin,
	// иначе - в файле.
	if g.flags.FileName == "" {
		g.stdinSearch()
	} else {
		g.fileSearch()
	}
}

// Метод, осуществляющий поиск совпадений в файле.
func (g *Grep) fileSearch() {
	// Мапа, в которую записываются индексы строк для вывода (строки-совпадения и строки до и после совпадения).
	indexesSet := make(map[int]struct{})
	// Количество строк в файле.
	fileContentLen := len(g.fileContent)

	for i, val := range g.fileContent {
		// Строка для сравнения.
		str := val
		// Паттерн, с которым сравнивается строка.
		pattern := g.flags.Pattern

		// Если при запуске программы указан флаг игнорирования регистра, то приводим строку
		// и паттерн к нижнему регистру.
		if g.flags.IgnoreCase {
			str = strings.ToLower(str)
			pattern = strings.ToLower(pattern)
		}

		// Поиск совпадений.
		switch {
		case g.flags.Invert:
			// Если при запуске указан флаг исключения, то результатом будут строки, которые не являются совпадением.
			if !g.compare(str, pattern) {
				indexesSet[i] = struct{}{}
				// Увеличиваем счетчик количества совпадений на единицу.
				g.matchesCount++
			}
		default:
			// Если строка является совпадением.
			if g.compare(str, pattern) {
				// Записываем индекс строки в результат.
				indexesSet[i] = struct{}{}

				// Записываем в результат индексы g.flags.After строк после совпадения.
				for j := i + 1; j < fileContentLen && j <= i+g.flags.After; j++ {
					indexesSet[j] = struct{}{}
				}

				// Записываем в результат индексы g.flags.Before строк до совпадения.
				for j := i - g.flags.Before; j < i; j++ {
					if j < 0 {
						continue
					}
					indexesSet[j] = struct{}{}
				}

				g.matchesCount++
			}
		}
	}

	// Если указан флаг вывода количества совпадений, то выводим их на экран и завершаем работу функции.
	if g.flags.Count {
		g.printMatchesCount()
		return
	}

	g.printStrings(g.createOrderedSet(indexesSet))
}

// Метод, сравнивающий строку и паттерн.
func (g Grep) compare(str, pattern string) bool {
	// Если указан флаг точного совпадения с паттерном, то возвращаем результат сравнения паттерна со строкой.
	if g.flags.Fixed {
		return str == pattern
	}
	// Иначе проверяем, содержит ли строка паттерн, и возвращаем результат.
	return strings.Contains(str, pattern)
}

// Из неупорядоченной мапы с индексами строк для вывода (indexesSet)
// создает отсортированный по возрастанию индексов слайс,
// который будет использоваться для вывода строк.
func (g Grep) createOrderedSet(indexesSet map[int]struct{}) []int {
	// Упорядоченное множество индексов строк для вывода.
	orderedSet := make([]int, len(indexesSet))
	i := 0

	// Проходимся по всем индексам в мапе и записываем их в слайс.
	for index := range indexesSet {
		orderedSet[i] = index
		i++
	}

	// Сортируем слайс по возрастанию.
	sort.Ints(orderedSet)

	return orderedSet
}

// Выводит на экран строки-совпадения, а также строки до и после совпадений, если указаны соответствующие флаги.
func (g Grep) printStrings(orderedSet []int) {
	// Проходимся по всем индексам строк для вывода.
	for _, index := range orderedSet {
		// Строка для вывода.
		str := g.fileContent[index]
		// Паттерн для сравнения.
		pattern := g.flags.Pattern

		// Если указан флаг игнорирования регистра, то приводим строку и паттерн к нижнему регистру.
		if g.flags.IgnoreCase {
			str = strings.ToLower(str)
			pattern = strings.ToLower(pattern)
		}

		// Строка с индексом выводимой строки.
		// Подсвечивается в консоли зеленым.
		// Непустая только в том случае, если указан флаг вывода номера строки.
		var lineNum string
		if g.flags.LineNum {
			lineNum = fmt.Sprintf("%s%d: ", greenColor, index+1)
		}

		// Если строка является строкой-совпадением, то она подсвечивается красным.
		if g.compare(str, pattern) {
			fmt.Print(lineNum, redColor, g.fileContent[index], "\n")
			continue
		}

		// Строки, которые не являются совпадениями, выводятся в обычном цвете.
		fmt.Print(lineNum, regularColor, g.fileContent[index], "\n")
	}

	// Устанавливаем в консоли обычный цвет шрифта.
	fmt.Print(regularColor)
}

// Выводит на экран количество строк-совпадений.
func (g Grep) printMatchesCount() {
	fmt.Println(g.matchesCount)
}

// Метод, осуществляющий поиск совпадений во время ввода данных в stdin.
func (g *Grep) stdinSearch() {
	reader := bufio.NewReader(os.Stdin)
	// Номер текущей строки.
	index := 1
	// Канал, в который записываются строки до совпадения.
	chBeforeData := newBeforeDataChannel(g.flags.Before)
	// Количество строк после совпадения, которые нужно вывести.
	afterCount := 0

	for {
		// s - введенная строка.
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Ошибка при чтении из stdin: %v", err)
		}
		// Слайс рун, в который будут скопированы руны из строки s, но без последних двух символов (10 и 13).
		sRunes := make([]rune, len([]rune(s))-2)
		copy(sRunes, []rune(s))
		// Присваиваем строке новое значение без лишних символов в конце строки.
		s = string(sRunes)

		if s != "" && s[0] == 4 && len(s) == 1 {
			break
		}

		// Строка для сравнения.
		str := s
		// Паттерн.
		pattern := g.flags.Pattern

		// Если указан флаг игнорирования регистра, то приводим строку и паттерн к нижнему регистру.
		if g.flags.IgnoreCase {
			str = strings.ToLower(str)
			pattern = strings.ToLower(pattern)
		}

		// Строка с индексом выводимой строки.
		// Подсвечивается в консоли зеленым.
		// Непустая только в том случае, если указан флаг вывода номера строки.
		var lineNum string
		if g.flags.LineNum {
			lineNum = fmt.Sprintf("%s%d: ", greenColor, index)
		}

		// Поиск совпадений.
		switch {
		case g.flags.Invert:
			if !g.compare(str, pattern) {
				// Если при запуске указан флаг исключения,
				// то результатом будут строки, которые не являются совпадением.
				// Вывод введенной строки на экран.
				fmt.Print(lineNum, regularColor, s, "\n")
				// Увеличиваем счетчик количества совпадений на единицу.
				g.matchesCount++
			}
		default:
			// Если строка является совпадением.
			if g.compare(str, pattern) {
				// Вывод на экран данных, которые были введены до совпадения.
				chBeforeData.PrintBeforeData(g.flags.LineNum)

				// Вывод на экран строки-совпадения, введенной раннее.
				fmt.Print(lineNum, redColor, s, "\n")
				fmt.Print(regularColor)

				// Количество строк, которые будут выведены после совпадения.
				// (Строки, введенные после совпадения, сразу выводятся на экран).
				afterCount = g.flags.After
				// Увеличиваем счетчик количества совпадений на единицу.
				g.matchesCount++
				// Увеличиваем номер текущей строки.
				index++
				continue
			}
		}

		// Если есть строки после совпадения для вывода, то выводим только что введенную строку на экран.
		if afterCount > 0 {
			fmt.Print(lineNum, regularColor, s, "\n")
			// Уменьшаем счетчик строк для вывода после совпадения на единицу.
			afterCount--
		} else if g.flags.Before > 0 {
			// Если строк для вывода после совпадения нет и если был указан ненулевой флаг для вывода
			// строк до совпадения, то записываем текущую строку в канал.
			chBeforeData.WriteBeforeData(index, s)
		}

		// Увеличиваем номер текущей строки.
		index++
	}

	// Если указан флаг вывода количества совпадений, то выводим их на экран.
	if g.flags.Count {
		g.printMatchesCount()
	}
}

func main() {
	grep := NewGrep()
	grep.Search()

	// Примеры работы программы:

	// 1) Поиск совпадений в файле
	// Содержание файла:
	// 1
	// 2
	// 2
	// 3
	// 2222
	// abcd
	// 12a
	// 324
	// абвгд
	//
	// 34
	// 2

	// Вывод программы с флагами -A 3 2 1:
	// 2
	// 2
	// 3
	// 2222
	// abcd
	// 12a
	// 324
	// абвгд
	//
	// 34
	// 2

	// Вывод программы с флагами -B 2 3 1:
	// 2
	// 2
	// 3
	// abcd
	// 12a
	// 324
	// абвгд
	//
	// 34

	// Вывод программы с флагами -C 1 3 1:
	// 2
	// 3
	// 2222
	// 12a
	// 324
	// абвгд
	//
	// 34
	// 2

	// Вывод программы с флагами -c 3 1:
	// 3

	// Вывод программы с флагами -C 1 -i Ab 1:
	// 2222
	// abcd
	// 12a

	// Вывод программы с флагами -v 2 1:
	// 1
	// 3
	// abcd
	// абвгд
	//
	// 34

	// Вывод программы с флагами -F 2 1:
	// 2
	// 2
	// 2

	// Вывод программы с флагами -C 1 -F -n 2 1:
	// 1: 1
	// 2: 2
	// 3: 2
	// 4: 3
	// 11: 34
	// 12: 2

	// 2) Поиск совпадений при вводе данных в stdin

	// Вывод программы с флагами -A 2 s:
	// 2
	// 3
	// sss
	// sss
	// 2
	// 2
	// 90
	// 90
	// 1
	// ^D

	// Вывод программы с флагами -B 2 -n 1:
	// 1
	// 1: 1
	// 23
	// 4
	// 1
	// 2: 23
	// 3: 4
	// 4: 1
	// 2
	// 1
	// 5: 2
	// 6: 1
	// ^D

	// Вывод программы с флагами -C 1 -n 1:
	// 1
	// 1: 1
	// 2
	// 2: 2
	// 3
	// 4
	// 5
	// 1
	// 5: 5
	// 6: 1
	// 6
	// 7: 6
	// 8
	// ^D

	// Вывод программы с флагами -c -n 1:
	// 2
	// 3
	// 1
	// 3: 1
	// 45
	// 12
	// 5: 12
	// 61
	// 6: 61
	// 2
	// ^D
	// 3

	// Вывод программы с флагами -i -n AAb:
	// a
	// aab
	// 2: aab
	// AAB
	// 3: AAB
	// f
	// ^D

	// Вывод программы с флагами -v -n 1:
	// 1
	// a
	// 2: a
	// b
	// 3: b
	// c
	// 4: c
	// 123
	// 2
	// 6: 2
	// ^D

	// Вывод программы с флагами -F -n 12:
	// 12
	// 1: 12
	// 123
	// 3124
	// ^D
}
