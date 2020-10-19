package stringutils

func Contains(s string, substrings ...string) (bool, bool, map[string]int) {
	if len(s) == 0 {
		return false, false, nil
	}

	if len(substrings) == 0 {
		return true, true, nil
	}

	m := make(map[rune][]string, len(substrings))

	for _, v := range substrings {
		if len(v) == 0 {
			continue
		}
		r := []rune(v)[0]
		val := v[1:]
		if len(v) == 1 {
			val = "nil"
		}
		if _, ok := m[r]; !ok {
			m[r] = []string{val}
		} else {
			tmp := m[r]
			tmp = append(tmp, val)
			m[r] = tmp
		}
	}

	var complete int
	m1 := make(map[string]int)

	for i, v := range s {
		if vals, ok := m[v]; ok {
			for _, val := range vals {
				if val == "nil" {
					complete++
					m1[string(v)] = i
					continue
				}
				if len(s)-i < len(val) {
					// fail
					continue
				}
				n := 1
				for _, r := range val {
					if r != []rune(s)[i+n] {
						break
					}
					if n == len(val) {
						m1[string(v)+val] = i
						complete++
						break
					}
					n++
				}
			}
		}
	}

	if complete > 0 {
		if complete == len(substrings) {
			return true, true, m1
		}
		return true, false, m1
	}
	return false, false, nil
}
