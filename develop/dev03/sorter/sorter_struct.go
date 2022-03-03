package sorter

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Структура содержит различные методы для сортировки файлов.
// Для сортировки используются функции из пакета sort.
type sorter struct {
	// Слайс строк файла.
	// Используется для обычной сортировки (не по столбцам).
	lines []string
	// Двумерный слайс, в котором содержатся слова из каждой строки.
	// Используется для сортировки по столбцам.
	matrixLines [][]string
	// Название исходного файла.
	sourceFileName string
	// Название файла, в который будет записан результат.
	resultFileName string
}

// Конструктор структуры sorter.
func newSorter(sourceFileName, resultFileName string) sorter {
	return sorter{
		lines:          make([]string, 0),
		matrixLines:    make([][]string, 0),
		sourceFileName: sourceFileName,
		resultFileName: resultFileName,
	}
}

// Производит чтение данных из исходного файла.
func (s *sorter) readFile() {
	// Открываем файл f.
	f, err := os.Open(s.sourceFileName)
	if err != nil {
		log.Fatalf("Не удалось открыть исходный файл. Ошибка: %v\n", err)
	}
	// После завершения работы функции файл закрывается.
	defer f.Close()

	// Построчно читаем данные из файла f.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s.lines = append(s.lines, scanner.Text())
	}
}

// Записывает отсортированные данные в файл результат.
// Если программа запускается со значением флага k > 0, то данные в файл записываются
// из двумерного слайса sorter.matrixLines.
func (s sorter) createResultFile(k int) {
	// Создаем файл результат f.
	f, err := os.Create(s.resultFileName)
	if err != nil {
		log.Fatalf("Не удалось создать файл. Ошибка: %v\n", err)
	}
	// После завершения работы функции файл закрывается.
	defer f.Close()

	// Запись отсортированных данных в файл.
	if k > 0 {
		for _, str := range s.matrixLines {
			f.WriteString(strings.Join(str, " ") + "\n")
		}
	} else {
		for _, str := range s.lines {
			f.WriteString(str + "\n")
		}
	}
}

// Обычная сортировка. Данный метод используется тогда, когда программа запускается без флага k.
// В функции реализованы также обратная сортировка, сортировка чисел и их комбинации.
// reverse - флаг обратной сортировки, num - флаг сортировки чисел.
func (s *sorter) simpleSort(reverse bool, num bool) {
	switch {
	case reverse:
		switch {
		case num:
			// Обратная сортировка чисел.
			// Работает по такому же принципу, что и обычная сортировка чисел,
			// но только данные сортируются в обратном порядке.
			// Пример сортировки:
			// 35				899
			// 46				46
			// 25				35
			// 899				25
			// 6		=>		6
			// ааа				вв
			// вв				ааа
			// qq				w wg
			// 65gg				qq
			// w wg				65gg
			sort.Slice(s.lines, func(i, j int) bool {
				si := s.lines[i]
				sj := s.lines[j]

				di, oki := toNumber(si)
				dj, okj := toNumber(sj)

				if oki && okj {
					return di > dj
				}
				if oki && !okj {
					return true
				}
				if !oki && okj {
					return false
				}
				return si > sj
			})

		default:
			// Обратная сортировка строк.
			// Пример сортировки:
			// нет специальной операции			нет специальной операции
			// 10								abcd
			// 1								5f f
			// 2						=>		2
			// abcd								10
			// 5f f								1
			sort.Slice(s.lines, func(i, j int) bool {
				return s.lines[i] > s.lines[j]
			})
		}

	default:
		switch {
		case num:
			// Сортировка чисел.
			// Работает таким образом, что строки, которые невозможно преобразовать в числа,
			// сортируются по обычному правилу (в данном случае сортировка по возрастанию).
			// А строки, которые преобразуются в числа, сортируются по правилу чисел.
			// В итоговом файле отсортированные числа будут находиться после отсортированных строк.
			// Пример сортировки:
			// 35				65gg
			// 46				qq
			// 25				w wg
			// 899				ааа
			// 6		=>		вв
			// ааа				6
			// вв				25
			// qq				35
			// 65gg				46
			// w wg				899
			sort.Slice(s.lines, func(i, j int) bool {
				si := s.lines[i]
				sj := s.lines[j]

				di, oki := toNumber(si)
				dj, okj := toNumber(sj)

				if oki && okj {
					return di < dj
				}
				if oki && !okj {
					return false
				}
				if !oki && okj {
					return true
				}
				return si < sj
			})

		default:
			// Сортировка строк.
			// Пример сортировки:
			// нет специальной операции			1
			// 10								10
			// 1								2
			// 2						=>		5f f
			// abcd								abcd
			// 5f f								нет специальной операции
			sort.Strings(s.lines)
		}
	}
}

// Создает двумерный слайс sorter.matrixLines.
func (s *sorter) createMatrix() {
	for _, val := range s.lines {
		s.matrixLines = append(s.matrixLines, strings.Split(val, " "))
	}
}

// Сортировка данных по столбцу. В аргументах k - номер столбца.
// В функции реализованы также обратная сортировка, сортировка чисел и их комбинации.
// reverse - флаг обратной сортировки, num - флаг сортировки чисел.
func (s *sorter) sortByColumn(k int, reverse, num bool) {
	switch {
	case reverse:
		switch {
		// Обратная сортировка чисел в столбце k.
		// Работает по такому же принципу, что и обычная сортировка чисел в столбце k,
		// но данные при этом сортируются в обратном порядке.
		// Пример обратной сортировки чисел в столбце 3:
		// She 12 -26378				to the 23
		// was 923 smiling 156			end. f -256
		// to the 23			=>		She 12 -26378
		// end. f -256					was 923 smiling 156
		// abcd							ff f
		// ff f							abcd
		case num:
			sort.Slice(s.matrixLines, func(i, j int) bool {
				iCount := len(s.matrixLines[i])
				jCount := len(s.matrixLines[j])

				var oki, okj bool
				var di, dj int

				if iCount >= k {
					di, oki = toNumber(s.matrixLines[i][k-1])
				}

				if jCount >= k {
					dj, okj = toNumber(s.matrixLines[j][k-1])
				}

				if oki && okj {
					return di > dj
				}
				if oki && !okj {
					return true
				}
				if !oki && okj {
					return false
				}
				return strings.Join(s.matrixLines[i], " ") > strings.Join(s.matrixLines[j], " ")
			})

		default:
			// Обратная сортировка слов в столбце k.
			// Работает по такому же принципу, что и обычная сортировка по столбцу, но только данные
			// сортируются в обратном порядке.
			// Пример обратной сортировки по столбцу 2:
			// Then,						going now. I'll
			// I'll be						I'll be
			// 2 3 0 a						come back when
			// going now. I'll		=>		it's all
			// come back when				2 3 0 a
			// 3d 3							3d 3
			// it's all						over.
			// over.						Then,
			sort.Slice(s.matrixLines, func(i, j int) bool {
				iCount := len(s.matrixLines[i])
				jCount := len(s.matrixLines[j])

				if iCount >= k && jCount >= k {
					return s.matrixLines[i][k-1] > s.matrixLines[j][k-1]
				}
				if iCount < k && jCount < k {
					return strings.Join(s.matrixLines[i], " ") > strings.Join(s.matrixLines[j], " ")
				}
				if iCount < k && jCount >= k {
					return false
				}
				return true
			})
		}

	default:
		switch {
		// Сортировка чисел в столбце k.
		// Слова в столбце k, которые можно преобразовать в числа, сортируются по правилу чисел.
		// А если слова в столбце преобразовать нельзя или если в строке столбца с номером k нет,
		// то эти строки сортируются по обычному правилу.
		// Строки, которые отсортированы по числу в столбце k, расположены после строк, отсортированных
		// по обычному правилу.
		// Пример сортировки чисел в столбце 3:
		// She 12 -26378				abcd
		// was 923 smiling 156			ff f
		// to the 23			=>		was 923 smiling 156
		// end. f -256					She 12 -26378
		// abcd							end. f -256
		// ff f							to the 23
		case num:
			sort.Slice(s.matrixLines, func(i, j int) bool {
				iCount := len(s.matrixLines[i])
				jCount := len(s.matrixLines[j])

				var oki, okj bool
				var di, dj int

				if iCount >= k {
					di, oki = toNumber(s.matrixLines[i][k-1])
				}

				if jCount >= k {
					dj, okj = toNumber(s.matrixLines[j][k-1])
				}

				if oki && okj {
					return di < dj
				}
				if oki && !okj {
					return false
				}
				if !oki && okj {
					return true
				}
				return strings.Join(s.matrixLines[i], " ") < strings.Join(s.matrixLines[j], " ")
			})

		default:
			// Сортировка слов в столбце k.
			// Сортировка работает таким образом, что все строки, которые содержат меньше слов, чем k,
			// сортируются по обычному правилу (в данном случае по возрастанию).
			// А строки, в которых есть столбец с номером k, сортируются по слову в этом столбце.
			// Все отсортированные строки по столбцу k будут находиться после отсортированных строк,
			// в которых столбца k нет.
			// Пример сортировки по столбцу 2:
			// Then,						Then,
			// I'll be						over.
			// 2 3 0 a						2 3 0 a
			// going now. I'll		=>		3d 3
			// come back when				it's all
			// 3d 3							come back when
			// it's all						I'll be
			// over.						going now. I'll
			sort.Slice(s.matrixLines, func(i, j int) bool {
				iCount := len(s.matrixLines[i])
				jCount := len(s.matrixLines[j])

				if iCount >= k && jCount >= k {
					return s.matrixLines[i][k-1] < s.matrixLines[j][k-1]
				}
				if iCount < k && jCount < k {
					return strings.Join(s.matrixLines[i], " ") < strings.Join(s.matrixLines[j], " ")
				}
				if iCount < k && jCount >= k {
					return true
				}
				return false
			})
		}
	}
}

// Удаляет повторяющиеся строки в слайсе, либо удаляет строки с повторяющимися столбцами k
// в двумерном слайсе.
// isColumn - флаг проверки повторяющихся слов в столбце.
// k - номер столбца.
func (s *sorter) createSet(isColumn bool, k int) {
	switch {
	case isColumn:
		// Удаление строк исходного файла с повторяющимся столбцом k.
		// Если строка не содержит столбца k, то считается, что столбец k в строке равен "".
		// Таким образом, в результат попадет только одна строка, состоящая из < k строк.
		// Пример сортировки без повторяющихся слов в столбце 2:
		// a				n
		// n				1 22 3
		// a a		=>		a a
		// c c c			s c d
		// a a				1234 g
		// s c d
		// 1 22 3
		// g g q
		// 1234 g

		// Мапа, в ключах которых будет храниться значение столбца k.
		// А в значении хранится массив слов строки.
		set := make(map[string][]string)
		// Двумерный слайс, куда записывается результат.
		linesSet := make([][]string, 0)

		for _, str := range s.matrixLines {
			if len(str) >= k {
				set[str[k-1]] = str
				continue
			}
			set[""] = str
		}

		for _, value := range set {
			linesSet = append(linesSet, value)
		}

		s.matrixLines = linesSet

	default:
		// Удаление повторяющихся строк.
		// Примеры сортировки без повторяющихся строк:
		// 1				1
		// fn				2
		// b b a d			3
		// 2		=>		a a
		// 2				b b a d
		// 3				fn
		// fn
		// a a

		// Мапа, в ключах которой будут храниться строки.
		set := make(map[string]struct{})
		// Слайс из неповторяющихся строк.
		linesSet := make([]string, 0)

		for _, str := range s.lines {
			set[str] = struct{}{}
		}

		for key := range set {
			linesSet = append(linesSet, key)
		}

		s.lines = linesSet
	}
}

// Преобразует строку в число.
// Если строка преобразуется в число, возвращает значение числа и true.
// Если строку невозможно преобразовать в число, то возвращается 0 и false.
func toNumber(s string) (int, bool) {
	d, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return d, true
}
