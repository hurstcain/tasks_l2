package dev02

//package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// UnpackString - функция, реализующая распаковку строки, поддерживает escape-последовательности.
// Примеры:
// `a4bc2d5e` -> `aaaabccddddde`;
// `qwe\\5` -> `qwe\\\\\`;
// `\8\9\0\\` -> `890\`;
// `g0j4` -> `jjjj` (Для функции также была реализована поддержка цифры 0. Если после символа стоит 0,
// то символ не добавляется в строку-результат).
// Если в функцию была передана некорректная исходная строка, то функция возвращает непустую ошибку.
func UnpackString(s string) (string, error) {
	// Слайс рун, содержащихся в строке s.
	var runes = []rune(s)
	// Количество рун в строке s.
	var runesCount = len(runes)
	// Если в строке первым элементом является цифра, то строка некорректная, возвращается ошибка.
	if runesCount > 0 && unicode.IsDigit(runes[0]) {
		return "", errors.New("invalid string")
	}
	// Указатель на структуру strings.Builder. В данную переменную будет записываться строка-результат.
	var result = new(strings.Builder)
	// Последняя прочитанная руна, используется, когда нужно добавить руну в строку-результат несколько раз.
	var lastRune rune

	for i := 0; i < runesCount; i++ {
		r := runes[i]

		// Если текущая руна r - число.
		if unicode.IsDigit(r) {
			// Если следующая руна тоже число, то строка s - некорректная, возвращается ошибка.
			if i+1 < runesCount && unicode.IsDigit(runes[i+1]) {
				return "", errors.New("invalid string")
			}
			// Преобразование руны в число.
			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			multiplyRune(result, lastRune, count)
			continue
		}

		// Если текущая руна r - обратная косая черта.
		if r == '\\' {
			// Увеличиваем счетчик на единицу.
			i++
			// Если косая черта - последний символ строки, то строка s - некорректная, возвращается ошибка.
			if i == runesCount {
				return "", errors.New("invalid string")
			}
			// Текущая руна теперь та, которая следует после обратной косой черты.
			r = runes[i]
		}

		// Если следующая руна - число, то меняем значение переменной lastRune и прерываем текущую итерацию цикла,
		// чтобы в результат не добавился лишний символ.
		if i+1 < runesCount && unicode.IsDigit(runes[i+1]) {
			lastRune = r
			continue
		}

		// Добавление текущей руны в результат.
		result.WriteRune(r)
	}

	return result.String(), nil
}

// Добавляет в строку-результат s руну r count раз.
func multiplyRune(s *strings.Builder, r rune, count int) {
	for i := 0; i < count; i++ {
		s.WriteRune(r)
	}
}
