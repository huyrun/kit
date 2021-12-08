package util

func RemoveDupString(items []string) []string {
	hash := make(map[string]struct{})
	for _, i := range items {
		hash[i] = struct{}{}
	}

	res := make([]string, len(hash))
	idx := 0
	for k := range hash {
		res[idx] = k
		idx++
	}

	return res
}

func RemoveAllChars(s string, chars []string) string {
	hash := make(map[string]struct{})
	for _, i := range chars {
		hash[i] = struct{}{}
	}

	var res string
	for _, c := range s {
		if _, ok := hash[string(c)]; !ok {
			res += string(c)
		}
	}

	return res
}

func GreaterOrEqualInt(arr []int, guard int) []int {
	res := make([]int, 0)
	for _, n := range arr {
		if n >= guard {
			res = append(res, n)
		}
	}

	return res
}
