package main

func Pair2020(in []int) (v [2]int, ok bool) {
	for _, i := range in {
		for _, j := range in {
			if i == j {
				continue
			}
			if i+j == 2020 {
				return [2]int{i, j}, true
			}
		}
	}
	return [2]int{}, false
}

func Triple2020(in []int) (v [3]int, ok bool) {
	for _, i := range in {
		for _, j := range in {
			for _, k := range in {
				if i == j || j == k || i == k {
					continue
				}
				if i+j+k == 2020 {
					return [3]int{i, j, k}, true
				}
			}
		}
	}
	return [3]int{}, false
}
