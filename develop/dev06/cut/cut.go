package cut

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// Cut - выводит определенные колонки строк, введенных раннее в stdin.
// columns - номера колонок для вывода.
// d - флаг, указывающий разделитель между колонками.
// s - флаг, указывающий выводить стоки без разделителей или нет.
func Cut(columns []int, d string, s bool) {
	matrix := readStrings(d)

	fmt.Printf("Результат:\n")
	printColumns(matrix, columns, s, d)
}

// Читает строки и разбивает их на колонки с помощью разделителя d.
// Возвращает данные, разделенные по строкам и колонкам и записанные в двумерный слайс.
func readStrings(d string) [][]string {
	matrix := make([][]string, 0)
	reader := bufio.NewReader(os.Stdin)

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
		// Если строка пустая, то выходим из цикла, то есть заканчиваем считывание новых строк.
		if s == "" {
			break
		}
		matrix = append(matrix, strings.Split(s, d))
	}

	return matrix
}

// Выводит на экран колонки, номера которых указаны в слайсе columns.
// Колонки разделяются разделителем d.
func printColumns(matrix [][]string, columns []int, s bool, d string) {
	// Количество колонок для вывода.
	columnsCount := len(columns)

	for _, line := range matrix {
		lineCount := len(line)
		// Если строка состоит из одной колонки, то проверяем, нужно ли выводить строки
		// без разделителя (флаг s). Если нужно, то выводим, если нет, то переходим на следующую итерацию цикла.
		if lineCount == 1 {
			if !s {
				fmt.Println(line[0])
			}
			continue
		}

		for i, column := range columns {
			// Если в строке у всех колонок номера меньше, чем номера колонок в columns, то на месте этой строки
			// выводится пустая строка, а цикл завершается.
			if column > lineCount-1 {
				fmt.Println()
				break
			}
			// Если колонка последняя или если колонка является последней для текущей строки или если
			// в строке больше нет колонок с номерами для вывода,
			// то добавляем в конец колонки символ перехода на новую строку и завершаем цикл.
			if i == columnsCount-1 || column == lineCount-1 || columns[i+1] >= lineCount {
				fmt.Printf("%s\n", line[column])
				break
			}
			fmt.Printf("%s%s", line[column], d)
		}
	}
}
