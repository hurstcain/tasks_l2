package dev01

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
)

// PrintCurrentTime - печатает на экран текущее время.
func PrintCurrentTime() {
	// Получаем текущее время с использованием ntp сервера.
	t, err := ntp.Time("ntp5.stratum2.ru")
	// В случае возникновения ошибки выводим содержание ошибки в os.Stderr и выходим из программы с кодом 13.
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(13)
	}
	// Если ошибки нет, то выводим отформатированное время.
	// Пример вывода:
	// Current time: 18:02:39.265 +03(MSK)
	fmt.Println(t.Format("Current time: 15:01:05.000 -07(MST)"))
}
