package common

func ContainsString(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, item string) []string {
	for idx, v := range slice {
		if v == item {
			return append(slice[:idx], slice[idx+1:]...)
		}
	}
	return slice
}
