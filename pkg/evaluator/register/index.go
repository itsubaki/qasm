package register

func Index(index, len int) int {
	if index < 0 {
		return len + index
	}

	return index
}
