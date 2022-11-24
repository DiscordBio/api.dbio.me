package middlewares

func Contains(array []string, value string) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}

func IndexOf(array []string, value string) int {
	for index, item := range array {
		if item == value {
			return index
		}
	}
	return -1
}

func Count(array []string) int {
	count := 0
	for range array {
		count++
	}
	return count
}

func Remove(array []string, value string) []string {
	index := IndexOf(array, value)
	if index == -1 {
		return array
	}
	return append(array[:index], array[index+1:]...)
}
