package dev04

import (
	"sort"
	"strings"
)

// CreateAnagramSet - создает множества анаграмм из слов в слайсе arr.
func CreateAnagramSet(arr []string) map[string][]string {
	anagramSet := make(map[string][]string)

	for _, s := range arr {
		// Приводим буквы в слове s к нижнему регистру.
		sLow := strings.ToLower(s)

		// Если ключ со значением sLow уже существует, то переходим на следующую итерацию цикла.
		if _, ok := anagramSet[sLow]; ok {
			continue
		}

		// Проверяем, есть ли множество анаграмм слова sLow или можно ли создать новое множество.
		if key, ok := searchAnagramForStr(sLow, anagramSet); ok {
			anagramSet[key] = append(anagramSet[key], sLow)
		}
	}

	// В данном цикле происходит удаление множеств с одним элементом
	// и сортировка множеств с несколькими элементами.
	for key, val := range anagramSet {
		if len(val) == 1 {
			delete(anagramSet, key)
		} else {
			sort.Strings(anagramSet[key])
		}
	}

	return anagramSet
}

// Функция проверяет, есть ли в множествах set множество анаграмм слова str.
// Если анаграмм слова нет, то функция возвращает это слово и флаг со значением true.
// В этом случае в множества будет добавлено множество новых анаграмм.
// Если множество анаграмм слова есть, то, при условии, что в множестве нет слова str,
// возвращается значение ключа множества и флаг true.
// В этом случае в множество анаграмм будет добавлено новое слово str.
func searchAnagramForStr(str string, set map[string][]string) (string, bool) {
	for key, val := range set {
		if areAnagrams(str, key) {
			if isElementInSet(str, val) {
				return "", false
			}
			return key, true
		}
	}

	return str, true
}

// Функция проверяет, являются ли слова a и b анаграммами.
// Если да, то возвращает true, если нет, то false.
func areAnagrams(a, b string) bool {
	// Преобразовываем слова a и b в слайсы рун
	aRunes := []rune(a)
	bRunes := []rune(b)

	// Сортируем по возрастанию руны в слайсах aRunes и bRunes.
	sort.Slice(aRunes, func(i, j int) bool {
		return aRunes[i] < aRunes[j]
	})
	sort.Slice(bRunes, func(i, j int) bool {
		return bRunes[i] < bRunes[j]
	})

	// Строки будут равны, если порядок элементов в отсортированных массивах aRunes и bRunes одинаковый.
	return string(aRunes) == string(bRunes)
}

// Функция проверяет, есть ли элемент со значением element в массиве set.
// Если есть, то возвращает true, если нет, то false.
func isElementInSet(element string, set []string) bool {
	for _, val := range set {
		if element == val {
			return true
		}
	}

	return false
}
