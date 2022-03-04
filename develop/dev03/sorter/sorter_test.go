package sorter

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSimpleSort(t *testing.T) {
	// Тестирование сортировки строк

	source := "test_files/simple_sort_line.txt"
	s := newSorter(source, "")
	expected := `1
10
2
25
300
300
300 a
Xian Le
Xian Le
ancient country
ancient country
flower
had no interest
strange
there was a famous`

	s.readFile()
	s.simpleSort(false, false)

	result := strings.Join(s.lines, "\n")
	assert.Equal(t, expected, result)

	// Тестирование обратной сортировки строк

	source = "test_files/simple_sort_line.txt"
	s = newSorter(source, "")
	expected = `there was a famous
strange
had no interest
flower
ancient country
ancient country
Xian Le
Xian Le
300 a
300
300
25
2
10
1`

	s.readFile()
	s.simpleSort(true, false)

	result = strings.Join(s.lines, "\n")
	assert.Equal(t, result, expected)

	// Тестирование сортировки чисел

	source = "test_files/simple_sort_numbers.txt"
	s = newSorter(source, "")
	expected = `300 a
ancient country
had no interest
-643
0
0
1
2
10
25
45
45
238
300
300`

	s.readFile()
	s.simpleSort(false, true)

	result = strings.Join(s.lines, "\n")
	assert.Equal(t, expected, result)

	// Тестирование обратной сортировки чисел

	source = "test_files/simple_sort_numbers.txt"
	s = newSorter(source, "")
	expected = `300
300
238
45
45
25
10
2
1
0
0
-643
had no interest
ancient country
300 a`

	s.readFile()
	s.simpleSort(true, true)

	result = strings.Join(s.lines, "\n")
	assert.Equal(t, expected, result)

	// Тестирование сортировки с выводом уникальных строк

	source = "test_files/simple_sort_line.txt"
	s = newSorter(source, "")
	expected = `1
10
2
25
300
300 a
Xian Le
ancient country
flower
had no interest
strange
there was a famous`

	s.readFile()
	s.createSet(false, 0)
	s.simpleSort(false, false)

	result = strings.Join(s.lines, "\n")
	assert.Equal(t, expected, result)
}

func TestSortByColumn(t *testing.T) {
	// Тестирование сортировки по столбцу

	source := "test_files/sort_by_column_line.txt"
	s := newSorter(source, "")
	expected := `earnest
that
This 10
speech and behaviour
и and
their glazed roofs
their glazed
can redeem yourself
can redeem yourself
Originally there had been nothing there`
	result := make([]string, 0)
	resultStr := ""

	s.readFile()
	s.createMatrix()
	s.sortByColumn(2, false, false)

	for _, v := range s.matrixLines {
		result = append(result, strings.Join(v, " "))
	}
	resultStr = strings.Join(result, "\n")
	assert.Equal(t, expected, resultStr)

	// Тестирование обратной сортировки по столбцу

	source = "test_files/sort_by_column_line.txt"
	s = newSorter(source, "")
	expected = `Originally there had been nothing there
can redeem yourself
can redeem yourself
their glazed roofs
their glazed
speech and behaviour
и and
This 10
that
earnest`
	result = make([]string, 0)

	s.readFile()
	s.createMatrix()
	s.sortByColumn(2, true, false)

	for _, v := range s.matrixLines {
		result = append(result, strings.Join(v, " "))
	}
	resultStr = strings.Join(result, "\n")
	assert.Equal(t, expected, resultStr)

	// Тестирование сортировки чисел по столбцу

	source = "test_files/sort_by_column_number.txt"
	s = newSorter(source, "")
	expected = `a A
aa aa
q
w
7 -10 7
1 0 9
1 2
9 6
0 200
a 345 0`
	result = make([]string, 0)

	s.readFile()
	s.createMatrix()
	s.sortByColumn(2, false, true)

	for _, v := range s.matrixLines {
		result = append(result, strings.Join(v, " "))
	}
	resultStr = strings.Join(result, "\n")
	assert.Equal(t, expected, resultStr)

	// Тестирование обратной сортировки чисел по столбцу

	source = "test_files/sort_by_column_number.txt"
	s = newSorter(source, "")
	expected = `a 345 0
0 200
9 6
1 2
1 0 9
7 -10 7
w
q
aa aa
a A`
	result = make([]string, 0)

	s.readFile()
	s.createMatrix()
	s.sortByColumn(2, true, true)

	for _, v := range s.matrixLines {
		result = append(result, strings.Join(v, " "))
	}
	resultStr = strings.Join(result, "\n")
	assert.Equal(t, expected, resultStr)

	// Тестирование сортировки по столбцу с выводом уникальных строк

	source = "test_files/sort_by_column_line.txt"
	s = newSorter(source, "")
	expected = `that
This 10
и and
their glazed
can redeem yourself
Originally there had been nothing there`
	result = make([]string, 0)

	s.readFile()
	s.createMatrix()
	s.createSet(true, 2)
	s.sortByColumn(2, false, false)

	for _, v := range s.matrixLines {
		result = append(result, strings.Join(v, " "))
	}
	resultStr = strings.Join(result, "\n")
	assert.Equal(t, expected, resultStr)
}
