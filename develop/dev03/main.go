package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hurstcain/tasks_l2/develop/dev03/sorter"
	"os"
	"strings"
)

func main() {
	var sourceFile, resultFile string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Название исходного файла: ")
	sourceFile, _ = reader.ReadString('\n')
	fmt.Print("Название файла, куда будет записан результат: ")
	resultFile, _ = reader.ReadString('\n')
	sourceFile = strings.TrimSpace(sourceFile)
	resultFile = strings.TrimSpace(resultFile)

	kFlag := flag.Int("k", 0, "Номер колонки для сортировки")
	nFlag := flag.Bool("n", false, "Сортировка по числовому значению")
	rFlag := flag.Bool("r", false, "Сортировка в обратном порядке")
	uFlag := flag.Bool("u", false, "Не выводить повторяющиеся строки.")
	flag.Parse()

	sorter.Sort(sourceFile, resultFile, *kFlag, *nFlag, *rFlag, *uFlag)
}
