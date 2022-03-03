package sorter

func Sort(source, result string, k int, n, r, u bool) {
	s := newSorter(source, result)
	s.ReadFile()

	switch {
	case k > 0:
		s.CreateMatrix()
		if u {
			s.CreateSet(true, k)
		}
		s.SortByColumn(k, r, n)
	case k == 0:
		if u {
			s.CreateSet(false, 0)
		}
		s.SimpleSort(r, n)
	}

	s.CreateResultFile(k)
}
