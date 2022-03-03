package sorter

import (
	"log"
)

// Sort - функция сортировки файла.
// source - исходный файл.
// result - файл, в который записывается результат.
// k - флаг сортировки по столбцу.
// n - сортировка чисел.
// r - обратная сортировка.
// u - вывод уникальных строк.
func Sort(source, result string, k int, n, r, u bool) {
	// Экземпляр структуры sorter.
	s := newSorter(source, result)
	s.readFile()

	switch {
	case k > 0:
		// Сортировка по столбцам.
		s.createMatrix()
		if u {
			s.createSet(true, k)
		}
		s.sortByColumn(k, r, n)

	case k == 0:
		// Обычная сортировка.
		if u {
			s.createSet(false, 0)
		}
		s.simpleSort(r, n)

	default:
		log.Fatalln("Некорректное значение флага k")
	}

	s.createResultFile(k)

	log.Printf("Сортировка строк из файла %s завершена\n", source)
	log.Printf("Результат записан в файл %s\n", result)
}
