package my_bridge

func Acquire(bitmap *[8192]byte) (int, bool) {
	for i := range *bitmap {
		if (*bitmap)[i] == 255 {
			continue
		}

		for j := 7; j >= 0; j-- {
			if (*bitmap)[i]&(1<<j) == 0 {
				(*bitmap)[i] |= 1 << j

				return 8*i + (7 - j), true
			}
		}

		panic("unreachable")
	}

	return -1, false
}

func Release(bitmap *[8192]byte, index int) bool {
	skip := index / 8

	// ignore
	if skip >= 8192 {
		return false
	}

	offset := 7 - index%8

	if (*bitmap)[skip]&(1<<offset) == 0 {
		return false
	}

	(*bitmap)[skip] ^= 1 << offset

	return true
}
