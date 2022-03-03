package main

import (
	"flag"
	"github.com/hurstcain/tasks_l2/develop/dev03/sorter"
)

func main() {
	sourceFile := "text.txt"
	resultFile := "new.txt"
	kFlag := flag.Int("k", 0, "Номер колонки для сортировки")
	nFlag := flag.Bool("n", false, "Сортировка по числовому значению")
	rFlag := flag.Bool("r", false, "Сортировка в обратном порядке")
	uFlag := flag.Bool("u", false, "Не выводить повторяющиеся строки.")
	flag.Parse()

	sorter.Sort(sourceFile, resultFile, *kFlag, *nFlag, *rFlag, *uFlag)
}
