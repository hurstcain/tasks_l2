package dev04

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAnagramSet(t *testing.T) {
	testData := []struct {
		arr      []string
		expected map[string][]string
	}{
		{
			arr:      []string{""},
			expected: make(map[string][]string),
		},
		{
			arr:      []string{"чулок", "рамка"},
			expected: make(map[string][]string),
		},
		{
			arr: []string{"яГуар", "ДупелЬ", "яруга", "ПУДЕЛЬ"},
			expected: map[string][]string{
				"ягуар":  {"ягуар", "яруга"},
				"дупель": {"дупель", "пудель"},
			},
		},
		{
			arr: []string{"пехота", "просев", "потеха", "подвес", "подсев", "совдеп"},
			expected: map[string][]string{
				"пехота": {"пехота", "потеха"},
				"подвес": {"подвес", "подсев", "совдеп"},
			},
		},
		{
			arr: []string{"собачка", "розетка", "басочка", "кератоз", "отрезка", "соусник", "косинус"},
			expected: map[string][]string{
				"собачка": {"басочка", "собачка"},
				"розетка": {"кератоз", "отрезка", "розетка"},
				"соусник": {"косинус", "соусник"},
			},
		},
		{
			arr: []string{"собачка", "розетка", "басочка", "кератоз", "отрезка", "собачка"},
			expected: map[string][]string{
				"собачка": {"басочка", "собачка"},
				"розетка": {"кератоз", "отрезка", "розетка"},
			},
		},
	}

	for _, data := range testData {
		result := CreateAnagramSet(data.arr)

		assert.Equal(t, data.expected, result)
	}
}

func TestSearchAnagramForStr(t *testing.T) {
	testData := []struct {
		str            string
		set            map[string][]string
		expectedString string
		expectedBool   bool
	}{
		{
			str: "амарант",
			set: map[string][]string{
				"плюшка":  {"плюшка", "шлюпка"},
				"маранта": {"маранта", "амарант"},
				"басочка": {"басочка", "собака"},
			},
			expectedString: "",
			expectedBool:   false,
		},
		{
			str: "пика",
			set: map[string][]string{
				"заруб": {"заруб", "арбуз"},
				"кипа":  {"кипа"},
			},
			expectedString: "кипа",
			expectedBool:   true,
		},
		{
			str: "солод",
			set: map[string][]string{
				"гроза": {"гроза", "розга"},
			},
			expectedString: "солод",
			expectedBool:   true,
		},
	}

	for _, data := range testData {
		resultString, resultBool := searchAnagramForStr(data.str, data.set)

		assert.Equal(t, data.expectedString, resultString)
		assert.Equal(t, data.expectedBool, resultBool)
	}
}

func TestAreAnagrams(t *testing.T) {
	testData := []struct {
		a        string
		b        string
		expected bool
	}{
		{
			a:        "агулка",
			b:        "калуга",
			expected: true,
		},
		{
			a:        "колесо",
			b:        "оселок",
			expected: true,
		},
		{
			a:        "лагерь",
			b:        "слиток",
			expected: false,
		},
		{
			a:        "павлин",
			b:        "спикер",
			expected: false,
		},
	}

	for _, data := range testData {
		result := areAnagrams(data.a, data.b)

		assert.Equal(t, data.expected, result)
	}
}

func TestIsElementInSet(t *testing.T) {
	testData := []struct {
		element  string
		set      []string
		expected bool
	}{
		{
			element:  "крона",
			set:      []string{"оркан", "крона", "норка"},
			expected: true,
		},
		{
			element:  "опрос",
			set:      []string{"просо", "сопор"},
			expected: false,
		},
	}

	for _, data := range testData {
		result := isElementInSet(data.element, data.set)

		assert.Equal(t, data.expected, result)
	}
}
