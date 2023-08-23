package utils

func IndexOf(slice []string, item string) float64 {
	for index, current := range slice {
		if current == item {
			return float64(index)
		}
	}
	return -1
}
