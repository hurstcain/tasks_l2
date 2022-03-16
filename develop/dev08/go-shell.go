package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// ExecuteCommands - исполняет shell команды.
func ExecuteCommands(commands []string) error {
	for _, val := range commands {
		// Слайс, состоящий из названия команды и ее параметров.
		// Сначала из строки удаляются лишние пробелы в конце и в начале строки.
		// Затем строка разделяется на элементы слайса с помощью разелителя пробела.
		// Первый элемент слайса - название команды, остальные - параметры команды.
		command := strings.Split(strings.TrimSpace(val), " ")
		// Если команда на исполнение - cd, то меняем текущую директорию с помощью os.Chdir.
		if command[0] == "cd" {
			err := os.Chdir(strings.Join(command[1:], " "))
			if err != nil {
				return err
			}
			continue
		}
		// Исполняет команду и возвращает результат исполнения команды.
		out, err := exec.Command(command[0], command[1:]...).Output()
		if err != nil {
			return err
		}
		// Выводим на экран результат исполнения команды.
		fmt.Println(string(out))
	}

	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// s - введенная строка.
		s, err := reader.ReadString('\n')
		if err != nil {
			// Если были нажаты ctrl+D, то ввод данных прекращается и программа завершает работу.
			if err == io.EOF {
				return
			}
			log.Fatalf("Error when reading from stdin: %v", err)
		}
		// Слайс рун, в который будут скопированы руны из строки s, но без символа переноса строки (10).
		sRunes := make([]rune, len([]rune(s))-1)
		copy(sRunes, []rune(s))
		// Присваиваем строке новое значение без лишних символов в конце строки.
		s = string(sRunes)

		// Если пользователь ввел \quit, то программа завершает работу.
		if s == `\quit` {
			return
		}

		// Слайс введенных команд, разделенных символом |.
		commands := strings.Split(s, "|")

		if err := ExecuteCommands(commands); err != nil {
			fmt.Printf("Error when executing a command: %v\n", err)
		}
	}

	// Примеры исполнения программы:
	// ps -F | echo ABCD abcd
	// UID          PID    PPID  C    SZ   RSS PSR STIME TTY          TIME CMD
	// pcell       2062    1625  0  2856  5804   0 11:19 pts/1    00:00:00 /bin/bash --rcfile /snap/goland/173/plugins/terminal/jediterm-bash.in -i
	// pcell       3527    2062  0 308937 17044  0 12:21 pts/1    00:00:00 go run go-shell.go
	// pcell       3553    3527  0 175761 1836   0 12:21 pts/1    00:00:00 /tmp/go-build627019461/b001/exe/go-shell
	// pcell       3577    3553  0  2884  3440   0 12:25 pts/1    00:00:00 ps -F
	//
	// ABCD abcd
	//
	// pwd
	// /home/pcell/Документы/develop/dev08
	//
	// cd /home/pcell/Документы | pwd
	// /home/pcell/Документы
	//
	// ps -p 1694 | kill 1694 | ps -p 1694
	//    PID TTY          TIME CMD
	//   1694 ?        00:00:00 fsnotifier
	//
	//
	// Error when executing a command: exit status 1
}
