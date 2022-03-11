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

type GrepFlags struct {
	Pattern    string
	FileName   string
	After      int
	Before     int
	Count      bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNum    bool
}

func NewGrepFlags() GrepFlags {
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

	return GrepFlags{
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

type Grep struct {
	fileContent  []string
	matchesCount int
	flags        GrepFlags
}

func NewGrep() Grep {
	flags := NewGrepFlags()

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

func (g Grep) Search() {
	if g.flags.FileName == "" {
		g.stdinSearch()
	} else {
		g.fileSearch()
	}
}

func (g *Grep) fileSearch() {
	indexesSet := make(map[int]struct{})
	fileContentLen := len(g.fileContent)

	for i, val := range g.fileContent {
		str := val
		pattern := g.flags.Pattern

		if g.flags.IgnoreCase {
			str = strings.ToLower(str)
			pattern = strings.ToLower(pattern)
		}

		switch {
		case g.flags.Invert:
			if !g.compare(str, pattern) {
				indexesSet[i] = struct{}{}
				g.matchesCount++
			}
		default:
			if g.compare(str, pattern) {
				indexesSet[i] = struct{}{}

				for j := i + 1; j < fileContentLen && j <= i+g.flags.After; j++ {
					indexesSet[j] = struct{}{}
				}

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

	if g.flags.Count {
		g.printMatchesCount()
		return
	}

	g.printStrings(g.createOrderedSet(indexesSet))
}

func (g Grep) compare(str, pattern string) bool {
	if g.flags.Fixed {
		return str == pattern
	}
	return strings.Contains(str, pattern)
}

func (g Grep) createOrderedSet(indexesSet map[int]struct{}) []int {
	orderedSet := make([]int, len(indexesSet))
	i := 0

	for index := range indexesSet {
		orderedSet[i] = index
		i++
	}

	sort.Ints(orderedSet)

	return orderedSet
}

func (g Grep) printStrings(orderedSet []int) {
	redColor := "\033[31m"
	greenColor := "\033[32m"
	regularColor := "\033[0m"

	for _, index := range orderedSet {
		str := g.fileContent[index]
		pattern := g.flags.Pattern

		if g.flags.IgnoreCase {
			str = strings.ToLower(str)
			pattern = strings.ToLower(pattern)
		}

		var lineNum string
		if g.flags.LineNum {
			lineNum = fmt.Sprintf("%s%d: ", greenColor, index)
		}

		if g.compare(str, pattern) {
			fmt.Print(lineNum, redColor, g.fileContent[index], "\n")
			continue
		}

		fmt.Print(lineNum, regularColor, g.fileContent[index], "\n")
	}
	fmt.Print(regularColor)
}

func (g Grep) printMatchesCount() {
	fmt.Println(g.matchesCount)
}

func (g *Grep) stdinSearch() {
}

func main() {
	grep := NewGrep()
	grep.Search()
}
