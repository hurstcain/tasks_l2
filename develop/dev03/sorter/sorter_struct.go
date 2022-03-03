package sorter

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type sorter struct {
	lines          []string
	matrixLines    [][]string
	sourceFileName string
	resultFileName string
}

func newSorter(sourceFileName, resultFileName string) sorter {
	return sorter{
		lines:          make([]string, 0),
		matrixLines:    make([][]string, 0),
		sourceFileName: sourceFileName,
		resultFileName: resultFileName,
	}
}

func (s *sorter) ReadFile() {
	f, err := os.Open(s.sourceFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s.lines = append(s.lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (s sorter) CreateResultFile(k int) {
	f, _ := os.Create(s.resultFileName)

	if k > 0 {
		for _, str := range s.matrixLines {
			f.WriteString(strings.Join(str, " ") + "\n")
		}
	} else {
		for _, str := range s.lines {
			f.WriteString(str + "\n")
		}
	}
}

func (s *sorter) SimpleSort(reverse bool, num bool) {
	switch {
	case reverse:
		switch {
		case num:
			sort.Slice(s.lines, func(i, j int) bool {
				si := s.lines[i]
				sj := s.lines[j]

				di, oki := toNumber(si)
				dj, okj := toNumber(sj)

				if oki && okj {
					return di > dj
				}
				if oki && !okj {
					return true
				}
				if !oki && okj {
					return false
				}
				return si > sj
			})
		default:
			sort.Slice(s.lines, func(i, j int) bool {
				return s.lines[i] > s.lines[j]
			})
		}
	default:
		switch {
		case num:
			sort.Slice(s.lines, func(i, j int) bool {
				si := s.lines[i]
				sj := s.lines[j]

				di, oki := toNumber(si)
				dj, okj := toNumber(sj)

				if oki && okj {
					return di < dj
				}
				if oki && !okj {
					return false
				}
				if !oki && okj {
					return true
				}
				return si < sj
			})
		default:
			sort.Strings(s.lines)
		}
	}
}

func (s *sorter) CreateMatrix() {
	for _, val := range s.lines {
		s.matrixLines = append(s.matrixLines, strings.Split(val, " "))
	}
}

func (s *sorter) SortByColumn(k int, reverse, num bool) {
	switch {
	case reverse:
		switch {
		case num:
			sort.Slice(s.matrixLines, func(i, j int) bool {
				iCount := len(s.matrixLines[i])
				jCount := len(s.matrixLines[j])

				var oki, okj bool
				var di, dj int

				if iCount >= k {
					di, oki = toNumber(s.matrixLines[i][k-1])
				}

				if jCount >= k {
					dj, okj = toNumber(s.matrixLines[j][k-1])
				}

				if oki && okj {
					return di > dj
				}
				if oki && !okj {
					return true
				}
				if !oki && okj {
					return false
				}
				return strings.Join(s.matrixLines[i], " ") > strings.Join(s.matrixLines[j], " ")
			})
		default:
			sort.Slice(s.matrixLines, func(i, j int) bool {
				iCount := len(s.matrixLines[i])
				jCount := len(s.matrixLines[j])

				if iCount >= k && jCount >= k {
					return s.matrixLines[i][k-1] > s.matrixLines[j][k-1]
				}
				if iCount < k && jCount < k {
					return strings.Join(s.matrixLines[i], " ") > strings.Join(s.matrixLines[j], " ")
				}
				if iCount < k && jCount >= k {
					return false
				}
				return true
			})
		}
	default:
		switch {
		case num:
			sort.Slice(s.matrixLines, func(i, j int) bool {
				iCount := len(s.matrixLines[i])
				jCount := len(s.matrixLines[j])

				var oki, okj bool
				var di, dj int

				if iCount >= k {
					di, oki = toNumber(s.matrixLines[i][k-1])
				}

				if jCount >= k {
					dj, okj = toNumber(s.matrixLines[j][k-1])
				}

				if oki && okj {
					return di < dj
				}
				if oki && !okj {
					return false
				}
				if !oki && okj {
					return true
				}
				return strings.Join(s.matrixLines[i], " ") < strings.Join(s.matrixLines[j], " ")
			})
		default:
			sort.Slice(s.matrixLines, func(i, j int) bool {
				iCount := len(s.matrixLines[i])
				jCount := len(s.matrixLines[j])

				if iCount >= k && jCount >= k {
					return s.matrixLines[i][k-1] < s.matrixLines[j][k-1]
				}
				if iCount < k && jCount < k {
					return strings.Join(s.matrixLines[i], " ") < strings.Join(s.matrixLines[j], " ")
				}
				if iCount < k && jCount >= k {
					return true
				}
				return false
			})
		}
	}
}

func (s *sorter) CreateSet(isColumn bool, k int) {
	switch {
	case isColumn:
		set := make(map[string][]string)
		linesSet := make([][]string, 0)

		for _, str := range s.matrixLines {
			if len(str) >= k {
				set[str[k-1]] = str
				continue
			}
			set[""] = str
		}

		for _, value := range set {
			linesSet = append(linesSet, value)
		}

		s.matrixLines = linesSet

	default:
		set := make(map[string]struct{})
		linesSet := make([]string, 0)

		for _, str := range s.lines {
			set[str] = struct{}{}
		}

		for key, _ := range set {
			linesSet = append(linesSet, key)
		}

		s.lines = linesSet
	}
}

func toNumber(s string) (int, bool) {
	d, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return d, true
}
