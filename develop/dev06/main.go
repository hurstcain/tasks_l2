package main

import (
	"flag"
	"github.com/hurstcain/tasks_l2/develop/dev06/cut"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// Номера колонок для вывода.
	columns := make([]int, 0)
	flag.Func("f", "Номера колонок для вывода", func(fValue string) error {
		s := strings.Split(fValue, ",")

		// Преобразуем строковые данные в числа.
		for _, val := range s {
			d, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			columns = append(columns, d-1)
		}

		// Сортируем введенные номера колонок по возрастанию для упорядоченного вывода колонок.
		sort.Ints(columns)
		return nil
	})
	dFlag := flag.String("d", "\t", "Разделитель")
	sFlag := flag.Bool("s", false, "Только строки с разделителем")
	flag.Parse()

	cut.Cut(columns, *dFlag, *sFlag)

	// Входные данные:
	// 245:789	4567	M:4540	Admin	01:10:1980
	// 535:763	4987	M:3476	Sales	11:04:1978
	// 123:645
	// 290:720	2390	M:2900	Sales	10:06:1990

	// Примеры вывода программы:

	// -f 1
	// 245:789
	// 535:763
	// 123:645
	// 290:720

	// -f 2,3,5
	// 4567    M:4540  01:10:1980
	// 4987    M:3476  11:04:1978
	// 123:645
	// 2390    M:2900  10:06:1990

	// -f 2,4,3 -s
	// 4567    M:4540  Admin
	// 4987    M:3476  Sales
	// 2390    M:2900  Sales

	// -f 2 -d :
	// 789     4567    M
	// 763     4987    M
	// 645
	// 720     2390    M

	// -f 2,3 -d : -s
	// 789     4567    M:4540  Admin   01
	// 763     4987    M:3476  Sales   11
	// 645
	// 720     2390    M:2900  Sales   10

	// -f 1,4 -d :
	// 245:10
	// 535:04
	// 123
	// 290:06

	// -f 4,5 -d :
	// 10:1980
	// 04:1978
	//
	// 06:1990
}
